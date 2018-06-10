# ROCD

ROCD Startup process:

* Load configuration
* Check hard prerequisites or exit:
    * Read Access to **keystore**
    * Write Access to **statestore**
* Check soft prerequisites:
    * Disk Space
        * Perform garbage collection
        * Or start in **read-only** mode
    * Network Connectivity
        * Or warn with connectivity information:
            * interface up/down
            * modem stats
            * dns status
            * basic routing inormation
    * Time Syncronization
        * Force time syncronization
        * Or warn with time syncronization stats
* Load Plugins:
    * **Updater**:
        * Root permissions
        * Device-specific
        * Applys software updates
        * Interface: `update(current, target)`
    * **Inspector**:
        * Root permissions
        * Device-specific
        * Inspects system when root permissions are required
* Check provisioning from keystore
    * Or start in provisioning mode
    * After provisioning, exit
* Downgrade permissions
* Start Main







