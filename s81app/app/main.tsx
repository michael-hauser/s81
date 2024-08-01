"use client";

import React, { useState, useEffect } from 'react';
import useWebSocket, { ReadyState } from 'react-use-websocket';
import styles from './page.module.scss';
import Weather from './components/Weather/Weather';
import Forecast from './components/Forecast/Forecast';
import Subway from './components/Subway/Subway';
import SubwaySouth from './components/Subway/SubwaySouth';
import { Direction, mapSubwayData, SubwayArrival } from './models/subwayData';

type WebSocketMessage = {
  key: 'weather-data' | 'subway-a' | 'subway-b' | 'subway-c';
  value: string;
};

const Main: React.FC = () => {
  const WS_URL = 'ws://localhost:8081/ws';
  const { lastJsonMessage, readyState } = useWebSocket(WS_URL, {
    share: true,
    shouldReconnect: () => true,
  });

  // State for different message types
  const [weather, setWeather] = useState<WeatherData | undefined>(undefined);
  const [subwayA, setSubwayA] = useState<SubwayArrival[]>([]);
  const [subwayB, setSubwayB] = useState<SubwayArrival[]>([]);
  const [subwayC, setSubwayC] = useState<SubwayArrival[]>([]);
  const [subwayData, setSubwayData] = useState<SubwayArrival[]>([]);

  // Effect to handle new messages
  useEffect(() => {
    if (lastJsonMessage) {
      const message: WebSocketMessage = lastJsonMessage as WebSocketMessage;
      console.log('Received message:', message);

      // Update the state based on the message key
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

  // Effect to update subwayData every 100ms
  useEffect(() => {
    const measure = 200;
    let counter = 1;
    const intervalId = setInterval(() => {
      const a = subwayA.map((arrival) => ({ ...arrival, arrivalSeconds: arrival.arrivalSeconds - measure * counter * 0.001 }));
      const b = subwayB.map((arrival) => ({ ...arrival, arrivalSeconds: arrival.arrivalSeconds - measure * counter * 0.001 }));
      const c = subwayC.map((arrival) => ({ ...arrival, arrivalSeconds: arrival.arrivalSeconds - measure * counter * 0.001 }));
      const updatedSubwayData = [...a, ...b, ...c].sort((a, b) => a.arrivalMinutes - b.arrivalMinutes);
      setSubwayData(updatedSubwayData);
      counter++;
    }, measure);

    // Clear interval on component unmount
    return () => clearInterval(intervalId);
  }, [subwayA, subwayB, subwayC]);

  return (
    <main className={styles.main}>
      <Weather data={weather} />
      <Forecast data={weather} />
      <Subway arrivals={subwayData} direction={Direction.North} />
      <Subway arrivals={subwayData} direction={Direction.South} />
    </main>
  );
};

export default Main;
