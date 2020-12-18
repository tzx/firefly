import io from 'socket.io-client';
import { createLogger, LogLevelString } from 'bunyan';
import { config } from '../lib/config';
import * as utils from '../lib/utils';
import { AssetTradeMessage, IApp2AppMessage, IApp2AppMessageListener } from '../lib/interfaces';

const log = createLogger({ name: 'clients/app2app.ts', level: utils.constants.LOG_LEVEL as LogLevelString });

let socket: SocketIOClient.Socket
let listeners: IApp2AppMessageListener[] = [];

export const init = async () => {
  establishSocketIOConnection();
};

function subscribeWithRetry() {
  log.info(`App2App subscription: ${config.app2app.destinations.kat}`)
  socket.emit('subscribe', [config.app2app.destinations.kat], (err: any, data: any) => {
    if (err) {
      log.error(`App2App subscription failure (retrying): ${err}`);
      setTimeout(subscribeWithRetry, utils.constants.SUBSCRIBE_RETRY_INTERVAL);
      return;
    }
    log.info(`App2App subscription succeeded: ${JSON.stringify(data)}`);
  });
}

const establishSocketIOConnection = () => {
  let error = false;
  socket = io.connect(config.app2app.socketIOEndpoint, {
    transportOptions: {
      polling: {
        extraHeaders: {
          Authorization: 'Basic ' + Buffer.from(`${config.appCredentials.user}` +
            `:${config.appCredentials.password}`).toString('base64')
        }
      }
    }
  }).on('connect', () => {
    if (error) {
      error = false;
      log.info('App2App messaging Socket IO connection restored');
    }
    subscribeWithRetry();
  }).on('connect_error', (err: Error) => {
    error = true;
    log.error(`App2App messaging Socket IO connection error. ${err.toString()}`);
  }).on('error', (err: Error) => {
    error = true;
    log.error(`App2app messaging Socket IO error. ${err.toString()}`);
  }).on('data', (app2appMessage: IApp2AppMessage) => {
    log.trace(`App2App message ${JSON.stringify(app2appMessage)}`);
    try {
      const content: AssetTradeMessage = JSON.parse(app2appMessage.content);
      for (const listener of listeners) {
        listener(app2appMessage.headers, content);
      }
    } catch (err) {
      log.error(`App2App message error ${err}`);
    }

  }) as SocketIOClient.Socket;
};

export const addListener = (listener: IApp2AppMessageListener) => {
  listeners.push(listener);
};

export const removeListener = (listener: IApp2AppMessageListener) => {
  listeners = listeners.filter(entry => entry != listener);
};

export const dispatchMessage = (to: string, content: string) => {
  socket.emit('produce', {
    headers: {
      from: config.app2app.destinations.kat,
      to
    },
    content
  });
};

export const reset = () => {
  if (socket) {
    log.info('App2App Socket IO connection reset');
    socket.close();
    establishSocketIOConnection();
  }
};