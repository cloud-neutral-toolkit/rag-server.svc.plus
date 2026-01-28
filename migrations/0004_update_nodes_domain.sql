-- 0004_update_nodes_domain.sql
-- Update legacy node addresses to new svc.plus domain

-- Update specific Tokyo node mapping
UPDATE nodes 
SET 
    address = 'ha-proxy-jp.svc.plus',
    server_name = 'ha-proxy-jp.svc.plus'
WHERE 
    address = 'tky-connector.onwalk.net';

-- Update any other nodes still using onwalk.net domain
UPDATE nodes 
SET 
    address = REPLACE(address, '.onwalk.net', '.svc.plus'),
    server_name = REPLACE(server_name, '.onwalk.net', '.svc.plus')
WHERE 
    address LIKE '%.onwalk.net%' OR server_name LIKE '%.onwalk.net%';
