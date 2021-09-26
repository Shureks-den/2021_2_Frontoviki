fiber = require('fiber')
expd = require('expirationd')
box.cfg{listen = 3301}
box.once("sessions", function()
    s = box.schema.space.create('sessions')
    s:format({
        {name = 'Value', type = 'string'},
        {name = 'UserId', type = 'unsigned'},
        {name = 'ExpiresAt', type = 'unsigned'}
    })
    s:create_index('primary', {
        type = 'hash',
        parts = {'Value'}
    })
    s:create_index('secondary', {
        type = 'tree',
        parts = {'UserId'}, unique = false
    })
end)
box.schema.user.passwd('pass')
function is_tuple_expired(args, tuple)
  if (tuple[3] < fiber.time()) then return true end
  return false
  end
expd.run_task('sessions', s.id, is_tuple_expired)
