# Castro

High performance Open Tibia automatic account creator written in **Go**

# How

Castro provides lua bindings using a pool of lua states. Each request gets a state from the pool. If there are no states available a new one is created and later saved on the pool