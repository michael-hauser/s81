import React, { createContext, useContext } from 'react';
import useWebSocket, { ReadyState } from 'react-use-websocket';

interface WebSocketProviderProps {
  children: React.ReactNode;
}

const WebSocketContext = createContext<{
  sendJsonMessage: (message: any) => void;
  lastJsonMessage: null | any;
  readyState: ReadyState;
}>({
  sendJsonMessage: (message: any) => {},
  lastJsonMessage: null,
  readyState: ReadyState.CONNECTING,
});

export const WebSocketProvider: React.FC<WebSocketProviderProps> = ({ children }) => {
  const WS_URL = 'ws://localhost:8081/ws';
  const { sendJsonMessage, lastJsonMessage, readyState } = useWebSocket(WS_URL, {
    share: false,
    shouldReconnect: () => true,
  });

  return (
    <WebSocketContext.Provider value={{ sendJsonMessage, lastJsonMessage, readyState }}>
      {children}
    </WebSocketContext.Provider>
  );
};

export const useWebSocketContext = () => useContext(WebSocketContext);
