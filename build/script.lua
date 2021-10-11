box.cfg {
    ------------------------
    -- Network configuration
    ------------------------

    -- The read/write data port number or URI
    -- Has no default value, so must be specified if
    -- connections will occur from remote clients
    -- that do not use “admin address”
    listen = '*:5555';
    -- listen = '*:3301';

    pid_file = "tarantool.pid";
    background = false;
    -- The server is considered to be a Tarantool replica
    -- it will try to connect to the master
    -- which replication_source specifies with a URI
    -- for example konstantin:secret_password@tarantool.org:3301
    -- by default username is "guest"
    -- replication_source="127.0.0.1:3102";

    -- The server will sleep for io_collect_interval seconds
    -- between iterations of the event loop
    io_collect_interval = nil;

    -- The size of the read-ahead buffer associated with a client connection
    readahead = 16320;

    ----------------------
    -- Memtx configuration
    ----------------------

    -- An absolute path to directory where snapshot (.snap) files are stored.
    -- If not specified, defaults to /var/lib/tarantool/INSTANCE
    -- memtx_dir = nil;

    -- How much memory Memtx engine allocates
    -- to actually store tuples, in bytes.
    memtx_memory = 128 * 1024 * 1024; -- 128Mb

    -- Size of the smallest allocation unit, in bytes.
    -- It can be tuned up if most of the tuples are not so small
    memtx_min_tuple_size = 16;

    -- Size of the largest allocation unit, in bytes.
    -- It can be tuned up if it is necessary to store large tuples
    memtx_max_tuple_size = 128 * 1024 * 1024; -- 128Mb

    ----------
    -- Logging
    ----------

    -- How verbose the logging is. There are six log verbosity classes:
    -- 1 – SYSERROR
    -- 2 – ERROR
    -- 3 – CRITICAL
    -- 4 – WARNING
    -- 5 – INFO
    -- 6 – VERBOSE
    -- 7 – DEBUG
    log_level = 5;

    -- By default, the log is sent to /var/log/tarantool/INSTANCE.log
    -- If logger is specified, the log is sent to the file named in the string
    log = "tarantool.log";
}

box.schema.user.create('proxy', {password = 'proxy_pass'})
box.schema.user.grant('proxy', 'read,write', 'universe')

s = box.schema.space.create('requests', {if_not_exists = true})
s:format({{name = 'id', type = 'unsigned'},{name = 'method', type = 'string'},
        {name = 'path', type = 'string'},{name = 'params', type = 'map'},
        {name = 'headers'},{name = 'cookies', type = 'map'},{name = 'body'}})    
s:create_index('primary', {type = 'hash', parts = {1, 'num'}})

s = box.schema.space.create('responses', {if_not_exists = true})
s:format({{name = 'id', type = 'unsigned'},{name = 'message', type = 'string'},
        {name = 'code', type = 'unsigned'}, {name = 'body', type = 'string'},
        {name = 'headers', type = 'map'}})
s:create_index('primary', {type = 'hash', parts = {1, 'num'}})
