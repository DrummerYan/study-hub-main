#!/bin/bash

echo "======================================"
echo "验证角色数据是否正确"
echo "======================================"

docker exec -i gva-mysql mysql -uroot -pAa@6447985 gva << 'EOF'

SELECT '所有角色及其父级关系：' AS info;
SELECT 
  a1.authority_id AS '角色ID',
  a1.authority_name AS '角色名',
  a1.parent_id AS '父级ID',
  COALESCE(a2.authority_name, '根节点') AS '父级名'
FROM sys_authorities a1
LEFT JOIN sys_authorities a2 ON a1.parent_id = a2.authority_id
ORDER BY a1.parent_id, a1.authority_id;

SELECT '' AS '';
SELECT '检查循环引用（应该为空）：' AS info;
SELECT 
  authority_id, 
  authority_name, 
  parent_id 
FROM sys_authorities 
WHERE parent_id = authority_id;

EOF

echo ""
echo "======================================"
echo "验证完成！"
echo "如果看到循环引用，请联系我进一步修复。"
echo "======================================"
