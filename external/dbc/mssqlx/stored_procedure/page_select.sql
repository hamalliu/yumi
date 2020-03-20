IF EXISTS
    (
    SELECT
        *
    FROM
        dbo.sysobjects
    WHERE
        id = object_id(N'[dbo].[page_select_v1]')
    AND
        OBJECTPROPERTY(id, N'IsProcedure') = 1
    )
    DROP PROCEDURE [dbo].[sys_Page_v2]
GO
CREATE PROCEDURE [dbo].[page_select_v1]
    @curPage int output,
    @total int output,

    @tableName NVARCHAR(255),
    @primaryKeyColName NVARCHAR(125),
    @selectJoin NVARCHAR(2000),
    @whereJoin NVARCHAR(2000),
    @where NVARCHAR(2000),
    @descOrder NVARCHAR(1000),
    @order NVARCHAR(1000),
    @fields NVARCHAR(2000),
    @page int,
    @line int
AS

BEGIN

--关掉提示信息，提高性能
SET NOCOUNT ON
SET ANSI_WARNINGS ON

--查询总条数
SET @totalsql =
        N'SELECT
             @total = COUNT(*)
         FROM
             @tableName
         @whereJoin
         @where;';

EXEC SP_EXECUTESQL @totalsql,
    N'@tableName NVARCHAR(255), @whereJoin NVARCHAR(2000), @where NVARCHAR(2000), @total int OUTPUT',
    @tableName, @whereJoin, @where, @total OUTPUT ;

--计算当前页数
SET @maxPage = ceiling(@total + 0.0/@line)
IF @page < 1
BEGIN
    SET @curPage = 1;
END
ELSE IF @page > @maxPage
BEGIN
   SET @curPage = @maxPage
END
ELSE
BEGIN
    SET @curPage = @page
END

--查询当页数据
SET @lastLine = @total-(@curPage-1)*@line

SET @dataSql =
    N'SELECT ' +
    N'   @fields ' +
    N'FROM ' +
    N'   @tableName ' +
    N'@selectJoin ' +
    N'WHERE ' +
    N'   @primaryKeyColName IN ' +
    N'       (' +
    N'       SELECT ' +
    N'           TOP @line ' +
    N'           ids.id ' +
    N'       FROM ' +
    N'           (' +
    N'           SELECT ' +
    N'           DISTINCT ' +
    N'           TOP @lastLine ' +
    N'               @primaryKeyColName AS id ' +
    N'           FROM ' +
    N'               @tableName ' +
    N'           @whereJoin ' +
    N'           @where ' +
    N'           @descOrder' +
    N'           ) AS ids ' +
    N'       @order' +
    N'       )';

EXEC SP_EXECUTESQL @dataSql,
    N'@fields NVARCHAR(2000), 
    @tableName NVARCHAR(255), 
    @selectJoin NVARCHAR(2000), 
    @primaryKeyColName NVARCHAR(125), 
    @line INT, 
    @lastLine INT, 
    @primaryKeyColName NVARCHAR(125), 
    @tableName NVARCHAR(255), 
    @whereJoin NVARCHAR(2000), 
    @where NVARCHAR(2000) 
    @descOrder NVARCHAR(1000) 
    @order NVARCHAR(1000)',
    @fields,
    @tableName,
    @selectJoin,
    @primaryKeyColName,
    @line,
    @lastLine,
    @primaryKeyColName,
    @tableName,
    @whereJoin,
    @where,
    @descOrder,
    @order

END

