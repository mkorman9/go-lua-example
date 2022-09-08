function handler(event)
    debug("[%s] event triggered by %s", event.timestamp, event.trigger)
end

events.add_listener(handler)
debug("registered event listener")
