package workflowengine

import (
	"fmt"
	"strconv"
	"strings"
	"time"
)

//操作流程实例============================================================================================================

func (m *Model) DeleteActor(actorId uint, proxyActorId uint) error {
	//如果还有待办未处理完，不能删除账号
	var (
		sql string

		err error
	)

	if proxyActorId == 0 {
		if count, err := m.GetTodoCount(actorId); err != nil {
			return err
		} else if count != 0 {
			return fmt.Errorf("该执行者还有未执行任务，请设定代理人")
		}
	}

	if err = m.MoveTask(actorId, proxyActorId, ""); err != nil {
		return err
	}

	m.db.Begin()
	sql = fmt.Sprintf("SELECT role_id FROM role_actors WHERE actor_id=%d", actorId)
	if rows, err := m.db.Raw(sql).Rows(); err != nil {
		m.db.Rollback()
		return err
	} else {
		for rows.Next() {
			roleId := 0
			rl := role{ID: uint(roleId)}
			if err = rows.Scan(&roleId); err != nil {
				m.db.Rollback()
				return err
			}
			if err = m.db.Related(&rl.Actors).Error; err != nil {
				m.db.Rollback()
				return err
			}
			m.actors[fmt.Sprintf("%d", roleId)] = rl.Actors
		}
	}

	sql = fmt.Sprintf("DELETE FROM role_actors WHERE actor_id=%d", actorId)
	if err = m.db.Exec(sql).Error; err != nil {
		m.db.Rollback()
		return err
	}
	m.db.Commit()

	return nil
}

func (m *Model) MoveTask(actorId uint, proxyActorId uint, wfCode string) error {
	var (
		err error
	)

	sql := fmt.Sprintf("UPDATE workflow_instances AS wfins "+
		"LEFT JOIN node_instances AS ndins ON ndins.workflow_instance_id=wfins.id "+
		"SET nd.actor_id=%d WHERE nd.actor_id=%d AND nd.status='%s'", proxyActorId, actorId, StatusNodeInstanceUndone)
	if wfCode != "" {
		sql = fmt.Sprintf("%s AND wfins.code='%s'", sql, wfCode)
	}
	if err = m.db.Exec(sql).Error; err != nil {
		return err
	}

	return nil
}

/**
 *发起流程
 *@mp:推动函数参数
 *@wsId:工作间id
 *@parWfinsId:父级流程实例id
 *@parWfinsId:父级流程实例id
 */
func (m *Model) Launch(wf string, parWfinsId uint, wsId uint, ctx string, rolePath string, actorId uint) error {
	var (
		outflow  string
		canNexts []string

		err error
	)

	//生成发起实例
	if parWfinsId != 0 {
		pwfins := WorkflowInstance{ID: parWfinsId}
		if err = m.db.Find(&pwfins).Error; err != nil {
			return err
		}

		if m.wfs[wf].ParentWfCode != pwfins.Code {
			return fmt.Errorf("流程：%s 父级流程不是流程：%s", m.wfs[wf].Name, pwfins.Name)
		}
	}
	wfins := WorkflowInstance{
		Year:          time.Now().Format("2006"),
		Name:          m.wfs[wf].Name,
		Code:          m.wfs[wf].Code,
		ParentWfinsID: parWfinsId,
		RolePath:      rolePath,
		ActorID:       actorId,
		UpdateDate:    time.Now(),
	}
	if err = m.db.Save(&wfins).Error; err != nil {
		return err
	}
	firstNodeIns := NodeInstance{
		Pahse:              m.nds[wf][m.wfs[wf].FirstNode].Pahse,
		WfCode:             m.nds[wf][m.wfs[wf].FirstNode].WfCode,
		NodeCode:           m.nds[wf][m.wfs[wf].FirstNode].Code,
		NodeName:           m.nds[wf][m.wfs[wf].FirstNode].Name,
		StartDate:          time.Now(),
		EndDate:            time.Now(),
		Deadline:           time.Now().AddDate(0, int(m.nds[wf][m.wfs[wf].FirstNode].ValidityPeriod), 0),
		WsName:             m.nds[wf][m.wfs[wf].FirstNode].WsName,
		RvwType:            m.nds[wf][m.wfs[wf].FirstNode].RvwType,
		Status:             StatusNodeInstanceStatusDone,
		WsId:               wsId,
		Context:            ctx,
		RolePath:           rolePath,
		ActorID:            actorId,
		UpdateDate:         time.Now(),
		WorkflowInstanceID: wfins.ID,
	}

	//获取流程实例
	if wfins, err = m.getWorkflowInstance(wfins.ID); err != nil {
		return err
	}

	//执行钩子函数
	if err = m.hook[m.nds[wf][m.wfs[wf].FirstNode].Hook](wfins); err != nil {
		return err
	}

	switch firstNodeIns.RvwType {
	case RvwtypeOrd, RvwtypeAssiOrd:
		if outflow, err = m.moveOrderNodeInstance(firstNodeIns); err != nil {
			return err
		}
	case RvwtypeGroup, RvwtypeAssiGroup:
		if outflow, err = m.moveGroupNodeInstance(firstNodeIns); err != nil {
			return err
		}
	default:
		if outflow, err = m.moveNodeInstance(firstNodeIns); err != nil {
			return err
		}
	}
	if outflow == "" {
		return nil
	}

	if m.nds[wf][firstNodeIns.NodeCode].OutflowMap[outflow] == nil {
		return fmt.Errorf("流程%s, 节点%s,不存在流向%s", wf, firstNodeIns.NodeCode, outflow)
	} else {
		canNexts = m.nds[wf][firstNodeIns.NodeCode].OutflowMap[outflow]
	}

	for i := range canNexts {
		if m.wfinsStat[canNexts[i]] != "" {
			switch m.wfinsStat[canNexts[i]] {
			case WfinsStatusTypePendingSet:
				if err = m.setWfinsStatus(wfins.ID, canNexts[i]); err != nil {
					return err
				}
			case WfinsStatusTypeSet:
				if err = m.pendingSetWfinsStatus(wfins.ID, canNexts[i]); err != nil {
					return err
				}
			}
			continue
		}

		if m.nds[wf][canNexts[i]] == nil {
			return fmt.Errorf("流程：%s， 节点：%s，下一步节点不存在：%s", wf, firstNodeIns.NodeCode, canNexts[i])
		}

		nxnode := m.nds[wf][canNexts[i]]

		//执行拦截器
		if ok, err := m.intercept[nxnode.Intercept](wfins); err != nil {
			return err
		} else if !ok {
			continue
		}

		if err = m.buildNodeInstance(wfins, *nxnode, firstNodeIns); err != nil {
			return err
		}

		//执行钩子函数
		if err = m.hook[nxnode.Hook](wfins); err != nil {
			return err
		}
	}

	if err = m.db.Save(&firstNodeIns).Error; err != nil {
		return err
	}

	return nil
}

func (m *Model) Move(wf string, wfinsId, nodeInsId, wsId uint, ctx string, rolePath string) error {
	var (
		nodeins  = NodeInstance{ID: nodeInsId}
		outflow  string
		canNexts []string
		wfins    WorkflowInstance

		err error
	)

	if err = m.db.Find(&nodeins).Error; err != nil {
		return err
	}
	nodeins.WsId = wsId
	nodeins.Context = ctx
	nodeins.RolePath = rolePath
	nodeins.UpdateDate = time.Now()

	//推动当前节点
	switch nodeins.RvwType {
	case RvwtypeOrd, RvwtypeAssiOrd:
		if outflow, err = m.moveOrderNodeInstance(nodeins); err != nil {
			return err
		}
	case RvwtypeGroup, RvwtypeAssiGroup:
		if outflow, err = m.moveGroupNodeInstance(nodeins); err != nil {
			return err
		}
	default:
		if outflow, err = m.moveNodeInstance(nodeins); err != nil {
			return err
		}
	}
	if outflow == "" {
		return nil
	}

	if m.nds[wf][nodeins.NodeCode].OutflowMap[outflow] == nil {
		return fmt.Errorf("流程%s, 节点%s,不存在流向%s", wf, nodeins.NodeCode, outflow)
	} else {
		canNexts = m.nds[wf][nodeins.NodeCode].OutflowMap[outflow]
	}

	//获取流程实例
	if wfins, err = m.getWorkflowInstance(wfinsId); err != nil {
		return err
	}

	for i := range canNexts {
		if m.wfinsStat[canNexts[i]] != "" {
			switch m.wfinsStat[canNexts[i]] {
			case WfinsStatusTypePendingSet:
				if err = m.setWfinsStatus(wfinsId, canNexts[i]); err != nil {
					return err
				}
			case WfinsStatusTypeSet:
				if err = m.pendingSetWfinsStatus(wfinsId, canNexts[i]); err != nil {
					return err
				}
			}
			continue
		}

		if m.nds[wf][canNexts[i]] == nil {
			return fmt.Errorf("流程：%s，节点：%s，下一步节点不存在：%s", wf, nodeins.NodeCode, canNexts[i])
		}

		nxnode := m.nds[wf][canNexts[i]]

		//执行拦截器
		if ok, err := m.intercept[nxnode.Intercept](wfins); err != nil {
			return err
		} else if !ok {
			continue
		}

		if err = m.buildNodeInstance(wfins, *nxnode, nodeins); err != nil {
			return err
		}

		//执行钩子函数
		if err = m.hook[nxnode.Hook](wfins); err != nil {
			return err
		}
	}

	//更新当前实例
	nodeins.EndDate = time.Now()
	if err = m.db.Save(&nodeins).Error; err != nil {
		return err
	}
	return nil
}

func (m *Model) existTodoNode(wfinsId uint) (bool, error) {
	var (
		count uint

		err error
	)

	if err = m.db.Model(&NodeInstance{}).Where("workflow_instance_id=?", wfinsId).Count(&count).Error; err != nil {
		return false, err
	}

	if count != 0 {
		return false, nil
	} else {
		return true, nil
	}
}

func (m *Model) setWfinsStatus(wfinsId uint, status string) error {
	var err error

	if err = m.db.Model(&WorkflowInstance{}).
		Where("id=?", wfinsId).Update("Status", status).Error; err != nil {
		return err
	}

	return nil
}

func (m *Model) pendingSetWfinsStatus(wfinsId uint, status string) error {
	var (
		ok bool

		err error
	)

	if ok, err = m.existTodoNode(wfinsId); err != nil {
		return err
	}
	if !ok {
		if err = m.db.Model(&WorkflowInstance{}).
			Where("id=?", wfinsId).Update("Status", status).Error; err != nil {
			return err
		}
	}

	return nil
}

func (m *Model) moveNodeInstance(nodeins NodeInstance) (string, error) {
	var (
		outflow string

		err error
	)

	if outflow, err = m.workspace[nodeins.WsName]([]uint{nodeins.WsId}); err != nil {
		return "", err
	}

	return outflow, nil
}

func (m *Model) moveOrderNodeInstance(nodeins NodeInstance) (string, error) {
	var (
		outflow string

		err error
	)

	//调用getoutflow函数
	if outflow, err = m.workspace[nodeins.WsName]([]uint{nodeins.WsId}); err != nil {
		return "", err
	}

	//生成同节点的下一个实例
	as := strings.Split(nodeins.NdinsActors, ",")
	if nodeins.NdinsSerialNum < uint(len(as)) && outflow == "" {
		actor, _ := strconv.Atoi(as[nodeins.NdinsSerialNum])
		nextNodeins := NodeInstance{
			NdinsSerialNum:     nodeins.NdinsSerialNum + 1,
			NdinsActors:        nodeins.NdinsActors,
			Pahse:              nodeins.Pahse,
			WfCode:             nodeins.WfCode,
			NodeCode:           nodeins.NodeCode,
			NodeName:           nodeins.NodeName,
			StartDate:          time.Now(),
			Deadline:           nodeins.Deadline,
			RvwType:            nodeins.RvwType,
			Status:             StatusNodeInstanceUndone,
			WsName:             nodeins.WsName,
			RolePath:           nodeins.RolePath,
			ActorID:            uint(actor),
			UpdateDate:         time.Now(),
			ParNodeInstanceID:  nodeins.ID,
			WorkflowInstanceID: nodeins.WorkflowInstanceID,
		}
		if err = m.db.Save(&nextNodeins).Error; err != nil {
			return "", err
		}
	}

	return outflow, nil
}

func (m *Model) moveGroupNodeInstance(nodeins NodeInstance) (string, error) {
	var (
		outflow string
		wsIds   []uint

		err error
	)

	//调用getoutflow函数
	sql := fmt.Sprintf("SELECT ws_id FROM node_instances WHERE workflow_instance_id=%d AND node_code='%s'",
		nodeins.WorkflowInstanceID, nodeins.NodeCode)
	if rows, err := m.db.Raw(sql).Rows(); err != nil {
		return "", err
	} else {
		for rows.Next() {
			var wsId uint
			if err = rows.Scan(wsId); err != nil {
				return "", err
			}
			wsIds = append(wsIds, wsId)
		}
	}
	if outflow, err = m.workspace[nodeins.WsName](wsIds); err != nil {
		return "", err
	}

	return outflow, nil
}

//如果该node已在流程实例中存在，那么将清除现有节点
func (m *Model) buildNodeInstance(wfins WorkflowInstance, node Node, prntNodeins NodeInstance) error {
	var (
		err error
	)

	switch prntNodeins.RvwType {
	case RvwtypeOrd:
		if err = m.buildOrderNodeInstance(wfins, node, prntNodeins); err != nil {
			return err
		}
	case RvwtypeGroup, "":
		if err = m.buildGroupNodeInstance(wfins, node, prntNodeins); err != nil {
			return err
		}
	case RvwtypeAssiOrd:
		if err = m.buildAssiOrdNodeInstance(wfins, node, prntNodeins); err != nil {
			return err
		}
	case RvwtypeAssiGroup:
		if err = m.buildAssiGroupNodeInstance(wfins, node, prntNodeins); err != nil {
			return err
		}
	default:
		return fmt.Errorf("流程：%s，节点：%s，不支持的节点类型%s", wfins.Code, node.Code, node.RvwType)
	}

	return nil
}

func (m *Model) buildOrderNodeInstance(wfins WorkflowInstance, node Node, prntNodeins NodeInstance) error {
	var (
		nxNodeins NodeInstance
		rolePath  string

		err error
	)

	if node.ReferenceNodes == "" {
		rolePath = prntNodeins.RolePath
	} else {
		rolePath = wfins.NodeInstancesMap[node.ReferenceNodes][0].RolePath
	}

	if node.RoleMap[rolePath] == nil {
		return fmt.Errorf("流程%s,节点%s,不存在rolepath:%s", node.WfCode, node.Code, rolePath)
	} else {
		rs := node.RoleMap[rolePath]
		for ri := range rs {
			rstr := fmt.Sprintf("%d", rs[ri])
			for ai := range m.actors[rstr] {
				if nxNodeins.ActorID == 0 {
					nxNodeins.ActorID = m.actors[rstr][ai].ID
				}
				if nxNodeins.NdinsActors == "" {
					nxNodeins.NdinsActors = fmt.Sprintf("%d", m.actors[rstr][ai].ID)
				} else {
					nxNodeins.NdinsActors = fmt.Sprintf("%s,%d", nxNodeins.NdinsActors, m.actors[rstr][ai].ID)
				}
			}
		}
	}

	nxNodeins.WfCode = node.WfCode
	nxNodeins.WorkflowInstanceID = prntNodeins.WorkflowInstanceID
	nxNodeins.WsName = node.WsName
	nxNodeins.Status = StatusNodeInstanceUndone
	nxNodeins.RvwType = node.RvwType
	nxNodeins.Deadline = time.Now().AddDate(0, int(node.ValidityPeriod), 0)
	nxNodeins.NodeName = node.Name
	nxNodeins.NodeCode = node.Code
	nxNodeins.Pahse = node.Pahse
	nxNodeins.NdinsSerialNum = 1
	nxNodeins.ParNodeInstanceID = prntNodeins.ID
	nxNodeins.StartDate = time.Now()
	nxNodeins.UpdateDate = time.Now()

	if err = m.db.Where("node_code='%s' AND workflow_instance_id=%d", node.Code, prntNodeins.WorkflowInstanceID).
		Delete(&NodeInstance{}).Error; err != nil {
		return err
	}
	if err = m.db.Save(&nxNodeins).Error; err != nil {
		return err
	}

	return nil
}

func (m *Model) buildGroupNodeInstance(wfins WorkflowInstance, node Node, prntNodeins NodeInstance) error {
	var (
		rolePath string
		count    uint

		err error
	)

	m.db.Begin()
	if err = m.db.Where("node_code='%s' AND workflow_instance_id=%d", node.Code, prntNodeins.WorkflowInstanceID).
		Delete(&NodeInstance{}).Error; err != nil {
		m.db.Rollback()
		return err
	}

	if node.ReferenceNodes == "" {
		rolePath = prntNodeins.RolePath
	} else {
		rolePath = wfins.NodeInstancesMap[node.ReferenceNodes][0].RolePath
	}

	if node.RoleMap[rolePath] == nil {
		m.db.Rollback()
		return fmt.Errorf("流程%s,节点%s,不存在rolepath:%s", node.WfCode, node.Code, rolePath)
	} else {
		rs := node.RoleMap[rolePath]
		for ri := range rs {
			rstr := fmt.Sprintf("%d", rs[ri])
			for ai := range m.actors[rstr] {
				count++
				nxNodeins := NodeInstance{}
				nxNodeins.ActorID = m.actors[rstr][ai].ID
				nxNodeins.WfCode = node.WfCode
				nxNodeins.WorkflowInstanceID = prntNodeins.WorkflowInstanceID
				nxNodeins.WsName = node.WsName
				nxNodeins.Status = StatusNodeInstanceUndone
				nxNodeins.RvwType = node.RvwType
				nxNodeins.Deadline = time.Now().AddDate(0, int(node.ValidityPeriod), 0)
				nxNodeins.NodeName = node.Name
				nxNodeins.NodeCode = node.Code
				nxNodeins.Pahse = node.Pahse
				nxNodeins.NdinsSerialNum = count
				nxNodeins.ParNodeInstanceID = prntNodeins.ID
				nxNodeins.StartDate = time.Now()
				nxNodeins.UpdateDate = time.Now()
				if err = m.db.Save(&nxNodeins).Error; err != nil {
					m.db.Rollback()
					return err
				}
			}
		}
	}
	m.db.Commit()

	return nil
}

func (m *Model) buildAssiOrdNodeInstance(wfins WorkflowInstance, node Node, prntNodeins NodeInstance) error {
	var (
		nxNodeins  NodeInstance
		assiActors string

		err error
	)

	if node.ReferenceNodes == "" {
		assiActors = prntNodeins.AssiActors
	} else {
		assiActors = wfins.NodeInstancesMap[node.ReferenceNodes][0].AssiActors
	}
	if assiActors == "" {
		return m.buildNodeInstance(wfins, *m.nds[node.WfCode][node.DefaultNode], prntNodeins)
	} else {
		as := strings.Split(assiActors, ",")
		for ai := range as {
			if nxNodeins.ActorID == 0 {
				actorId, _ := strconv.Atoi(as[ai])
				nxNodeins.ActorID = uint(actorId)
			}
			if nxNodeins.NdinsActors == "" {
				nxNodeins.NdinsActors = fmt.Sprintf("%s", as[ai])
			} else {
				nxNodeins.NdinsActors = fmt.Sprintf("%s,%s", nxNodeins.NdinsActors, as[ai])
			}
		}
	}

	nxNodeins.WfCode = node.WfCode
	nxNodeins.WorkflowInstanceID = prntNodeins.WorkflowInstanceID
	nxNodeins.WsName = node.WsName
	nxNodeins.Status = StatusNodeInstanceUndone
	nxNodeins.RvwType = node.RvwType
	nxNodeins.Deadline = time.Now().AddDate(0, int(node.ValidityPeriod), 0)
	nxNodeins.NodeName = node.Name
	nxNodeins.NodeCode = node.Code
	nxNodeins.Pahse = node.Pahse
	nxNodeins.NdinsSerialNum = 1
	nxNodeins.ParNodeInstanceID = prntNodeins.ID
	nxNodeins.StartDate = time.Now()
	nxNodeins.UpdateDate = time.Now()

	if err = m.db.Where("node_code='%s' AND workflow_instance_id=%d", node.Code, prntNodeins.WorkflowInstanceID).
		Delete(&NodeInstance{}).Error; err != nil {
		return err
	}
	if err = m.db.Save(&nxNodeins).Error; err != nil {
		return err
	}

	return nil
}

func (m *Model) buildAssiGroupNodeInstance(wfins WorkflowInstance, node Node, prntNodeins NodeInstance) error {
	var (
		assiActors string
		count      uint

		err error
	)

	m.db.Begin()
	if err = m.db.Where("node_code='%s' AND workflow_instance_id=%d", node.Code, prntNodeins.WorkflowInstanceID).
		Delete(&NodeInstance{}).Error; err != nil {
		m.db.Rollback()
		return err
	}

	if node.ReferenceNodes == "" {
		assiActors = prntNodeins.AssiActors
	} else {
		assiActors = wfins.NodeInstancesMap[node.ReferenceNodes][0].AssiActors
	}
	if assiActors == "" {
		if err = m.buildNodeInstance(wfins, *m.nds[node.WfCode][node.DefaultNode], prntNodeins); err != nil {
			m.db.Rollback()
		}
		return err
	} else {
		as := strings.Split(assiActors, ",")
		for ai := range as {
			count++
			nxNodeins := NodeInstance{}

			actorId, _ := strconv.Atoi(as[ai])
			nxNodeins.ActorID = uint(actorId)
			nxNodeins.WfCode = node.WfCode
			nxNodeins.WorkflowInstanceID = prntNodeins.WorkflowInstanceID
			nxNodeins.WsName = node.WsName
			nxNodeins.Status = StatusNodeInstanceUndone
			nxNodeins.RvwType = node.RvwType
			nxNodeins.Deadline = time.Now().AddDate(0, int(node.ValidityPeriod), 0)
			nxNodeins.NodeName = node.Name
			nxNodeins.NodeCode = node.Code
			nxNodeins.Pahse = node.Pahse
			nxNodeins.NdinsSerialNum = count
			nxNodeins.ParNodeInstanceID = prntNodeins.ID
			nxNodeins.StartDate = time.Now()
			nxNodeins.UpdateDate = time.Now()
			if err = m.db.Save(&nxNodeins).Error; err != nil {
				m.db.Rollback()
				return err
			}
		}
	}
	m.db.Commit()

	return nil
}

//todo
func (m *Model) canceledWfins() {

}

func (m *Model) getWorkflowInstance(wfinsId uint) (WorkflowInstance, error) {
	var (
		err error

		wfins = WorkflowInstance{ID: wfinsId}
	)

	if err = m.db.Find(&wfins).Related(&wfins.NodeInstances).Related(&wfins.SubWfinss).Error; err != nil {
		return wfins, err
	}

	wfins.NodeInstancesMap = make(map[string] /*code*/ []NodeInstance)
	wfins.SubWfinssMap = make(map[string] /*code*/ []WorkflowInstance)
	for ni := range wfins.NodeInstances {
		wfins.NodeInstancesMap[wfins.NodeInstances[ni].NodeCode] = append(wfins.NodeInstancesMap[wfins.NodeInstances[ni].NodeCode], wfins.NodeInstances[ni])
	}

	for swi := range wfins.SubWfinss {
		if wfins.SubWfinss[swi], err = m.getWorkflowInstance(wfins.SubWfinss[swi].ID); err != nil {
			return wfins.SubWfinss[swi], err
		}

		wfins.SubWfinssMap[wfins.SubWfinss[swi].Code] = append(wfins.SubWfinssMap[wfins.NodeInstances[swi].NodeCode], wfins.SubWfinss[swi])
	}

	return wfins, nil
}

//查询流程实例============================================================================================================

/**
 *获取单个流程实例
 *@wfId:流程实例id
 *@phase:阶段名称
 */
func (m *Model) GetWorkflowInstance(wfinsId uint, phase string) (WorkflowInstance, error) {
	var (
		wfins = WorkflowInstance{ID: wfinsId}

		err error
	)

	if err = m.db.Find(&wfins).Error; err != nil {
		return wfins, err
	}

	if phase != "" {
		if err = m.db.Where("phase = ?", phase).Find(&wfins.NodeInstances).Error; err != nil {
			return wfins, err
		}
	} else {
		if err = m.db.Model(&wfins).Related(&wfins.NodeInstances).Error; err != nil {
			return wfins, err
		}
	}

	return wfins, err
}

/**
 *获取流程发起者流程列表
 *@st, et:开始时间，截止时间
 *@actor:发起人
 *@wfinsstat:流程实例状态
 *@offset, line:分页
 */
func (m *Model) GetWorkflowInstanceListByInitiaor(
	sd, ed string, actorId uint, wfinsstat string, offset, line uint) ([]WorkflowInstance, int, error) {
	var (
		wfinss []WorkflowInstance
		total  int

		err error
	)

	sql := fmt.Sprintf("SELECT * FROM workflow_instances "+
		"WHERE actor_id=%d AND status='%s'", actorId, wfinsstat)
	if sd != "" {
		sql = fmt.Sprintf("%s AND start_date > '%s'", sql, sd)
	}
	if sd != "" {
		sql = fmt.Sprintf("%s AND start_date < '%s'", sql, ed)
	}
	if line != 0 {
		sql = fmt.Sprintf("%s LIMIT %d, %d", sql, offset, line)
	}
	if err = m.db.Raw(sql).Scan(&wfinss).Error; err != nil {
		return wfinss, total, err
	}

	if err = m.db.Raw(sql).Count(&total).Error; err != nil {
		return wfinss, total, err
	}

	return wfinss, total, nil
}

func (m *Model) GetWorkflowInstanceListByActor(
	sd, ed string, actorId uint, wfinsstat string, offset, line uint) ([]WorkflowInstance, int, error) {
	var (
		wfinss []WorkflowInstance
		total  int

		err error
	)

	sql := fmt.Sprintf("SELECT DISTINCT(wfins.*) FROM workflow_instances AS wfins "+
		"LEFT JOIN node_instances AS ndins ON wfins.id=ndins.workflow_instance_id AND ndins.actor_id=%d "+
		"WHERE wfins.status='%s'", actorId, wfinsstat)
	if sd != "" {
		sql = fmt.Sprintf("%s AND wfins.start_date > '%s'", sql, sd)
	}
	if sd != "" {
		sql = fmt.Sprintf("%s AND wfins.start_date < '%s'", sql, ed)
	}
	if line != 0 {
		sql = fmt.Sprintf("%s LIMIT %d, %d", sql, offset, line)
	}
	if err = m.db.Raw(sql).Scan(&wfinss).Error; err != nil {
		return wfinss, total, err
	}

	if err = m.db.Raw(sql).Count(&total).Error; err != nil {
		return wfinss, total, err
	}

	return wfinss, total, nil
}

func (m *Model) GetWorkflowInstanceListByReviewer(
	sd, ed string, actorId uint, wfinsstat string, offset, line uint) ([]WorkflowInstance, int, error) {
	var (
		wfinss []WorkflowInstance
		total  int

		err error
	)
	sql := fmt.Sprintf("SELECT DISTINCT(wfins.*) FROM workflow_instances AS wfins "+
		"LEFT JOIN node_instances AS ndins ON wfins.id=ndins.workflow_instance_id AND ndins.actor_id=%d "+
		"WHERE wfins.status='%s' AND wfins.actor_id<>'%d'", actorId, wfinsstat, actorId)
	if sd != "" {
		sql = fmt.Sprintf("%s AND wfins.start_date > '%s'", sql, sd)
	}
	if sd != "" {
		sql = fmt.Sprintf("%s AND wfins.start_date < '%s'", sql, ed)
	}
	if line != 0 {
		sql = fmt.Sprintf("%s LIMIT %d, %d", sql, offset, line)
	}
	if err = m.db.Raw(sql).Scan(&wfinss).Error; err != nil {
		return wfinss, total, err
	}

	if err = m.db.Raw(sql).Count(&total).Error; err != nil {
		return wfinss, total, err
	}

	return wfinss, total, nil
}

func (m *Model) GetTodo(actorId uint, offset, line uint) ([]WorkflowInstance, int, error) {
	var (
		wfinss []WorkflowInstance
		total  int

		err error
	)
	sql := fmt.Sprintf("SELECT wfins.* FROM workflow_instances AS wfins "+
		"LEFT JOIN node_instances AS ndins ON wfins.id=ndins.workflow_instance_id AND ndins.status='%s' AND ndins.actor_id=%d ",
		StatusNodeInstanceUndone, actorId)
	if line != 0 {
		sql = fmt.Sprintf("%s LIMIT %d, %d", sql, offset, line)
	}
	if err = m.db.Raw(sql).Scan(&wfinss).Error; err != nil {
		return wfinss, total, err
	}

	if err = m.db.Raw(sql).Count(&total).Error; err != nil {
		return wfinss, total, err
	}

	return wfinss, total, nil
}

func (m *Model) GetTodoCount(actor uint) (uint, error) {

	return 0, nil
}
