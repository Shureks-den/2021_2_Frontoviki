FROM tarantool/tarantool:1.10.2

WORKDIR /tarantool

COPY internal/database/tarantool_sandbox/start.lua start.lua
COPY internal/database/tarantool_sandbox/expirationd.lua expirationd.lua
COPY internal/database/tarantool_sandbox/checks.lua checks.lua

EXPOSE 3302

ENTRYPOINT [ "tarantool", "start.lua" ]