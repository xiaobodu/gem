
import sys
import argparse

import config
import signal_handler
import console
import plugins
import auth

import gem
import gem.runite
import gem.auth

from plugin_loader import PluginLoader
from service_listeners import ServiceListeners

version_string = "Gem v0.9: Opal"

plugin_path = ["content/plugins"]

# Create argparser
parser = argparse.ArgumentParser(description=version_string)
parser.add_argument('--console', action='store_true', help='launch the interactive console')
parser.add_argument('--plugin-path', action='append', dest='plugin_path', help='append to the plugin search path')

logger = gem.syslog.Module("pymain")
args = parser.parse_args()

def main():
    logger.Notice("Starting {0}".format(version_string))

    try:
        gem.runite.context = gem.runite.Context()
        gem.runite.context.Unpack(config.game_data['data_file'], config.game_data['index_files'])

        plugin_loader = PluginLoader(plugin_path)
        plugin_loader.load()

        # init service listeners
        # inserts an engine.Start hook to launch listeners
        listeners = ServiceListeners()
        listeners.register()

        # start the engine
        engine = gem.Engine()
        engine.Start()
        signal_handler.setup_exit_handler(engine.Stop)

        logger.Info("Finished engine initialization")
    except Exception as e:
        logger.Critical("Startup failed: {0}".format(e))

    # enter interactive console if --console flag is set
    if args.console:
        interactive_console()

    logger.Notice("Press Control-D to toggle the interactive console")
    while True:
        line = sys.stdin.readline()
        if not line: # readline will return "" on EOF
            interactive_console()

def interactive_console():
    logger.Notice("Transferring control to interactive console")
    gem.syslog.BeginRedirect()
    console.interact()
    gem.syslog.EndRedirect()
    logger.Info("Exited interactive console")

if __name__ == "__main__":
    main()
