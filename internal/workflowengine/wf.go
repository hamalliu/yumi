package workflowengine

import (
	"fmt"
	"strconv"
	"strings"
)

//操作流程定义============================================================================================================

func (m *Model) AddWorkflow(name string, codePrefix string) (uint, error) {
	var (
		wf      Workflow
		maxCode int

		err error
	)

	if name == "" ||
		codePrefix == "" {
		return 0, fmt.Errorf("name, codePrefix不能为空")
	}

	m.wfMux.Lock()

	wf.Name = name
	if maxCode, err = m.getMaxWfCode(); err != nil {
		return 0, err
	}
	wf.Code = codePrefix + fmt.Sprintf("_%d", maxCode+1)
	wf.Code = WorkflowStatusMaking

	if err = m.db.Save(&wf).Error; err != nil {
		return 0, fmt.Errorf("添加流程定义失败，%s", err.Error())
	}

	m.wfs[wf.Code] = &wf

	m.wfMux.Unlock()

	return wf.ID, nil
}

func (m *Model) UpdateWorkflow(name string, codePrefix string, id uint) error {
	var (
		wf      = Workflow{ID: id}
		updates = make(map[string]interface{})

		err error
	)

	if err = m.db.Find(&wf).Error; err != nil {
		return err
	}
	if codePrefix != "" {
		m.wfs[wf.Code].Code = codePrefix + strings.Split(wf.Code, "_")[1]
		updates["code"] = m.wfs[wf.Code].Code
	}
	if name != "" {
		m.wfs[wf.Code].Name = name
		updates["name"] = name
	}

	if len(updates) != 0 {
		if err = m.db.Model(&wf).Updates(updates).Error; err != nil {
			return err
		}
	}

	return nil
}

func (m *Model) ReleaseWorkflow(id uint) error {
	var (
		err error
	)

	if err = m.isValidWorkflow(id); err != nil {
		return err
	}

	if err = m.db.Model(&Workflow{ID: id}).Update("Status", WorkflowStatusRelease).Error; err != nil {
		return err
	}

	return nil
}

func (m *Model) DeleteWorkflow(wfId uint) error {
	var (
		wf    = Workflow{ID: wfId}
		count int

		err error
	)

	if err = m.db.Find(&wf).Error; err != nil {
		return err
	}

	if err = m.db.Model(&NodeInstance{}).
		Where("wf_code = ? AND status = ?", wf.Code, StatusNodeInstanceUndone).Count(&count).Error; err != nil {
		return nil
	}

	if count != 0 {
		return fmt.Errorf("该流程还有流程实例没有完结，不能删除。")
	}
	if err = m.db.Delete(wf).Error; err != nil {
		return err
	}
	delete(m.wfs, wf.Code)
	delete(m.nds, wf.Code)

	return nil
}

func (m *Model) AddNode(node Node, isFirstNode bool) (uint, error) {
	var (
		maxCode int
		wf      Workflow

		err error
	)

	node.ID = 0

	if node.WfCode == "" {
		return 0, fmt.Errorf("流程code不能为空")
	}

	wf.ID = node.WorkflowID
	if err = m.db.Find(&wf).Error; err != nil {
		return 0, err
	}
	if wf.Status == WorkflowStatusRelease {
		if err = m.isValidNode(node); err != nil {
			return 0, err
		}
	}

	m.wfs[node.WfCode].nodeMux.Lock()

	if maxCode, err = m.getMaxNodeCode(node.WorkflowID); err != nil {
		return 0, err
	}
	node.Code = fmt.Sprintf("N_%d", maxCode+1)

	m.db.Begin()
	if err = m.db.Save(&node).Error; err != nil {
		m.db.Rollback()
		return 0, fmt.Errorf("添加节点失败，%s", err.Error())
	}

	m.nds[node.WfCode][node.Code] = &node

	if err = m.db.Model(&Workflow{}).Where("code=?", node.WfCode).Update("first_node", node.Code).
		Error; err != nil {
		m.db.Rollback()
		return 0, err
	}
	m.db.Commit()
	m.wfs[node.WfCode].nodeMux.Unlock()

	return node.ID, nil
}

func (m *Model) UpdateNode(node Node) error {
	var (
		oldNode = Node{ID: node.ID}
		wf      Workflow

		err error
	)

	if node.ID == 0 {
		return fmt.Errorf("node的id不能为0")
	}
	wf.ID = node.WorkflowID
	if err = m.db.Find(&wf).Error; err != nil {
		return err
	}
	if wf.Status == WorkflowStatusRelease {
		if err = m.isValidNode(node); err != nil {
			return err
		}
	}

	if err = m.db.Find(&oldNode).Error; err != nil {
		return err
	}
	node.Code = oldNode.Code

	if err = m.db.Save(&node).Error; err != nil {
		return err
	}
	m.nds[node.WfCode][node.Code] = &node

	return nil
}

func (m *Model) DeleteNode(nodeId uint) error {
	var (
		node    = Node{ID: nodeId}
		wf      Workflow
		fdnodes string

		err error
	)

	if err = m.db.Find(&node).Error; err != nil {
		return err
	}
	wf.ID = node.WorkflowID
	if err = m.db.Find(&wf).Error; err != nil {
		return err
	}
	if wf.Status == WorkflowStatusRelease {
		if fdnodes, err = m.checkFlowDirectionNode(node.Code, node.WfCode); err != nil {
			return err
		}
	}

	if fdnodes != "" {
		return fmt.Errorf("无流向节点才能删除，流向%s节点的节点有%s", node.Code, fdnodes)
	}

	if err = m.db.Delete(&node).Error; err != nil {
		return err
	}
	delete(m.nds[node.WfCode], node.Code)

	return nil
}

func (m *Model) DeleteRole(roleId uint) error {
	//如果流程定义中还有该角色，不能删除
	for _, wv := range m.wfs {
		for ni := range wv.Nodes {
			for _, rs := range wv.Nodes[ni].RoleMap {
				for ri := range rs {
					if rs[ri] == roleId {
						return fmt.Errorf("流程:%s,节点:%s 存在该角色，请修改流程", wv.Code, wv.Nodes[ni].Code)
					}
				}
			}
		}
	}

	if err := m.db.Delete(&role{ID: roleId}).Error; err != nil {
		return err
	}
	delete(m.actors, fmt.Sprintf("%d", roleId))

	return nil
}

func (m *Model) DeleteActorOfRole(actorId uint, roleId uint) error {
	var (
		role role
	)
	sql := fmt.Sprintf("DELETE FROM role_actors WHERE role_id=%d AND actor_id=%d", roleId, actorId)
	if err := m.db.Exec(sql).Error; err != nil {
		return err
	}

	if err := m.db.Related(&role.Actors).Error; err != nil {
		return err
	}
	m.actors[fmt.Sprintf("%d", roleId)] = role.Actors

	return nil
}

func (m *Model) AddWorkflowInstanceStatus(t string, v string) (uint, error) {
	var (
		wfinsSt WorkflowInstanceStatus

		err error
	)

	wfinsSt.Type = t
	wfinsSt.Vaule = v
	if err = m.db.Save(&wfinsSt).Error; err != nil {
		return 0, err
	}

	return wfinsSt.ID, nil
}

func (m *Model) UpdateWorkflowInstanceStatus(id uint, t string, v string) error {
	var (
		updates = make(map[string]interface{})

		err error
	)

	if t != "" {
		updates["type"] = t
	}

	if v != "" {
		updates["value"] = v
	}

	if err = m.db.Model(&WorkflowInstanceStatus{ID: id}).Updates(updates).Error; err != nil {
		return err
	}

	return nil
}

func (m *Model) GetWorkflowInstanceStatus() error {
	var (
		wfinsSts []WorkflowInstanceStatus

		err error
	)

	if err = m.db.Find(&wfinsSts).Error; err != nil {
		return err
	}

	return nil
}

func (m *Model) DeleteWorkflowInstanceStatus(id uint) error {
	var (
		err error
	)

	if err = m.db.Delete(&WorkflowInstanceStatus{ID: id}).Error; err != nil {
		return err
	}

	return nil
}

func (m *Model) getMaxWfCode() (int, error) {
	var (
		maxCode int
		wf      Workflow

		err error
	)

	if err = m.db.Raw("SELECT wf FROM workflows ORDER BY id DESC").Scan(&wf).Error; err != nil {
		return maxCode, err
	}

	if maxCode, err = strconv.Atoi(strings.Split(wf.Code, "_")[1]); err != nil {
		return maxCode, err
	}

	return maxCode, nil
}

func (m *Model) getMaxNodeCode(wfId uint) (int, error) {
	var (
		maxCode int
		node    Node

		err error
	)

	if err = m.db.Raw("SELECT wf FROM nodes WHERE workflow_id=? ORDER BY id DESC", wfId).Scan(&node).Error; err != nil {
		return maxCode, err
	}

	if maxCode, err = strconv.Atoi(node.Code[1:]); err != nil {
		return maxCode, err
	}

	return maxCode, nil
}

func (m *Model) checkFlowDirectionNode(ncode, wfcode string) (string, error) {
	var (
		nodes   []Node
		fdnodes string

		err error
	)

	if err = m.db.Model(&Node{}).Where("wf_code = ?", wfcode).Scan(&nodes).Error; err != nil {
		return "", err
	}

	for i := range nodes {
		ofs := strings.Split(nodes[i].Outflow, ";")
		for ofi := range ofs {
			cns := strings.Split(strings.Split(ofs[ofi], ":")[1], ",")
			for cni := range cns {
				if ncode == cns[cni] {
					if fdnodes == "" {
						fdnodes = nodes[i].Code
					} else {
						fdnodes = fmt.Sprintf("%s,%s", fdnodes, nodes[i].Code)
					}
					continue
				}
			}
		}
	}

	return fdnodes, nil
}

func (m *Model) isValidWorkflow(wfId uint) error {
	var (
		wf Workflow

		err error
	)

	if err = m.db.Find(&wf).Related(&wf.Nodes).Error; err != nil {
		return err
	}

	for ni := range wf.Nodes {
		if err = m.isValidNode(wf.Nodes[ni]); err != nil {
			return err
		}
	}

	if wf.Name == "" {
		return fmt.Errorf("流程%s，流程名称不能为空", wf.Code)
	}
	if wf.FirstNode == "" {
		return fmt.Errorf("流程%s，未找到第一个节点", wf.Code)
	}

	wf.Status = WorkflowStatusRelease
	if err = m.db.Save(&wf).Error; err != nil {
		return err
	}

	return nil
}

func (m *Model) isValidNode(node Node) error {
	if node.WfCode == "" {
		return fmt.Errorf("节点%s，没有流程编码", node.Code)
	}

	if node.Name == "" {
		return fmt.Errorf("节点%s，没有节点名称", node.Code)
	}

	if node.WsName == "" {
		return fmt.Errorf("节点%s，没有工作间名称", node.Code)
	}

	if !isValidRoles(node.Roles) {
		return fmt.Errorf("节点%s，roles格式错误，格式：com>dep>role:1,2,3;com>dep>role:1,2,3", node.Code)
	}

	if !isValidOutflow(node.Outflow) {
		return fmt.Errorf("节点%s，outflow格式错误，格式：outflow:nodecode,nodecode;outflow:nodecode,nodecode", node.Code)
	}

	return nil
}

func isValidRoles(roles string) bool {
	if roles == "" {
		return false
	}
	roles = strings.Replace(roles, " ", "", -1)

	rs := strings.Split(roles, ";")
	for ri := range rs {
		if strings.Index(rs[ri], ":") == -1 {
			return false
		}
		if len(strings.Split(rs[ri], ":")) != 2 {
			return false
		}
		role := strings.Split(rs[ri], ":")[1]
		if role == "" {
			return false
		}
		for _, v := range strings.Split(role, ",") {
			if _, err := strconv.Atoi(v); err != nil {
				return false
			}
		}

	}

	return true
}

func isValidOutflow(outflow string) bool {
	if outflow == "" {
		return false
	}
	outflow = strings.Replace(outflow, " ", "", -1)

	os := strings.Split(outflow, ";")
	for oi := range os {
		if strings.Index(os[oi], ":") == -1 {
			return false
		}
		if len(strings.Split(os[oi], ":")) != 2 {
			return false
		}
		role := strings.Split(os[oi], ":")[1]
		if role == "" {
			return false
		}
		for _, v := range strings.Split(role, ",") {
			if len(v) < 4 {
				return false
			}
			if v[0] != 'N' || v[1] != '_' {
				return false
			}
			if _, err := strconv.Atoi(v[2:]); err != nil {
				return false
			}
		}

	}

	return true
}

//查询流程定义============================================================================================================
func (m *Model) GetWorkflowList(codePrefix string, offset, line uint) ([]Workflow, int, error) {
	var (
		wfs   []Workflow
		total int

		err error
	)

	if err = m.db.Find(&wfs).Offset(offset).Limit(line).Error; err != nil {
		return nil, 0, err
	}
	if err = m.db.Model(&Workflow{}).Count(&total).Error; err != nil {
		return nil, 0, err
	}

	return wfs, total, nil
}

func (m *Model) GetWorkflow(wfcode string) (Workflow, error) {
	var (
		wf = Workflow{Code: wfcode}

		err error
	)

	if err = m.db.Find(&wf).Related(&wf.Nodes).Error; err != nil {
		return wf, err
	}

	return wf, nil
}
