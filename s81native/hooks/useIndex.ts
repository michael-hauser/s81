import { useState, useEffect } from 'react';
import useWebSocket from 'react-use-websocket';
import { mapSubwayData, SubwayArrival } from '@/models/subwayData';

type WebSocketMessage = {
  key: 'weather-data' | 'subway-a' | 'subway-b' | 'subway-c';
  value: string;
};

export function useIndex() {
  const WS_URL = process.env.EXPO_PUBLIC_WS_URL;
  if (!WS_URL) {
    throw new Error('Missing process.env.EXPO_PUBLIC_WS_URL');
  }

  const { lastJsonMessage } = useWebSocket(WS_URL, {
    share: true,
    shouldReconnect: () => true,
  });

  const [weather, setWeather] = useState<WeatherData | undefined>(undefined);
  const [subwayA, setSubwayA] = useState<SubwayArrival[]>([]);
  const [subwayB, setSubwayB] = useState<SubwayArrival[]>([]);
  const [subwayC, setSubwayC] = useState<SubwayArrival[]>([]);

  useEffect(() => {
    if (lastJsonMessage) {
      const message: WebSocketMessage = lastJsonMessage as WebSocketMessage;

      switch (message.key) {
        case 'weather-data':
          setWeather(JSON.parse(message.value));
          break;
        case 'subway-a':
          setSubwayA(mapSubwayData(JSON.parse(message.value)));
          break;
        case 'subway-b':
          setSubwayB(mapSubwayData(JSON.parse(message.value)));
          break;
        case 'subway-c':
          setSubwayC(mapSubwayData(JSON.parse(message.value)));
          break;
        default:
          console.error('Unknown message key:', message.key);
      }
    }
  }, [lastJsonMessage]);

  const subwayData = [...subwayA, ...subwayB, ...subwayC].sort((a, b) => a.arrivalMinutes - b.arrivalMinutes);

  return { weather, subwayData };
}