# API接口

## 获取动态策略列表

### Request

- Method: **GET**
- URL:  ```/api/v1/policy/tag/lists?page={page}&pageSize={pageSize}&policy_name={策略名称}&name={标签}```

- <details>

  <summary>Body</summary>

  ``` json

  ```

</details>

### Response

- <details>

  <summary>Body</summary>

  ``` json
  {
    "msg": {
      "total": 1,
      "page": 1,
      "page_size": 25,
      "items": [
        {
          "id": "61404ff5bb4e33000671e522-61404ff4bb4e33000671e520",
          "name": [
            "WWW-*",
            "Administrator*"
          ],
          "ruleindex": 1,
          "policy_tag": [
            "system_defense"
          ],
          "policy_name": "123"
        }
      ]
    },
    "success": true
  }
  ```

  - 字段说明
    - ```total: 总条数```
    - ```page: 当前页数```
    - ```page_size: 当前页行数```
    - ```items: 列表详情```
      - ```id```
      - ```name: []string 标签数组```
      - ```ruleindex: 规则排名```
      - ```policy_tag: []string 策略配置```
      - ```policy_name: 策略名称```

</details>
