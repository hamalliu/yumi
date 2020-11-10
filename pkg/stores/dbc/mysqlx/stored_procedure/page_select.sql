DELIMITER $$
#分页查询存储过程
CREATE PROCEDURE page_select(IN p_cloumns varchar(2000), p_table varchar(255), p_where varchar(4000),
                            p_order varchar(500), p_index int, p_size int, OUT p_record_count int, p_cur_index int, p_count int)
begin
    declare v_sqlcounts varchar(4000);
    declare v_sqlselect varchar(4000);
    #查询总条数
    set v_sqlcounts = concat('SELECT COUNT(*) INTO @recordcount FROM ', p_table, p_where);
    set @sqlcounts = v_sqlcounts;
    prepare stmt from @sqlcounts;
    execute stmt;
    deallocate prepare stmt;
    set p_record_count = @recordcount;
    #根据总记录条数计算总页数
    set p_count = ceiling((p_record_count + 0.0) / p_size);
    if p_index < 1 then
        set p_cur_index = 1;
    elseif p_index > p_count then
        set p_cur_index = p_count;
    else
        set p_cur_index = p_index;
    end if;
    #查询当前页内容
    set v_sqlselect = concat('SELECT ', p_cloumns, 'FROM', p_table, p_where,
                             ifnull(p_order, ''), 'LIMIT ', (p_index-1)*p_size,',',p_size);
    set @sqlselect = v_sqlselect;
    prepare stmtselect from @sqlselect;
    execute stmtselect;
    deallocate prepare stmtselect;
end $$

DELIMITER ;