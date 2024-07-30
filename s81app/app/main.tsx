"use client";

import React, { useState, useEffect } from 'react';
import useWebSocket, { ReadyState } from 'react-use-websocket';
import styles from './page.module.scss';
import Weather from './components/Weather/Weather';
import Forecast from './components/Forecast/Forecast';
import Subway from './components/Subway/SubwayNorth';
import SubwaySouth from './components/Subway/SubwaySouth';

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
  const [subwayA, setSubwayA] = useState<SubwayData | undefined>(undefined);
  const [subwayB, setSubwayB] = useState<SubwayData | undefined>(undefined);
  const [subwayC, setSubwayC] = useState<SubwayData | undefined>(undefined);

  // Effect to handle new messages
  useEffect(() => {
    if (lastJsonMessage) {
      // Type the message
      const message: WebSocketMessage = lastJsonMessage as WebSocketMessage;

      console.log('Received message:', message);

      // Update the state based on the message key
      switch (message.key) {
        case 'weather-data':
          setWeather(JSON.parse(message.value));
          console.log('Weather data:', JSON.parse(message.value));
          break;
        case 'subway-a':
          setSubwayA(JSON.parse(message.value));
          break;
        case 'subway-b':
          setSubwayB(JSON.parse(message.value));
          break;
        case 'subway-c':
          setSubwayC(JSON.parse(message.value));
          break;
        default:
          console.error('Unknown message key:', message.key);
      }
    }
  }, [lastJsonMessage]);


  const toNYCTime = (time: number | undefined) => {
    if (!time) {
      return '-';
    }

    // Convert Unix timestamp to JavaScript Date object
    const date = new Date(time * 1000); // Unix timestamp is in seconds, so multiply by 1000 to convert to milliseconds

    // Define options for formatting date and time in NYC time
    const options: Intl.DateTimeFormatOptions = {
      timeZone: 'America/New_York',
      year: 'numeric',
      month: 'long',
      day: 'numeric',
      hour: 'numeric',
      minute: 'numeric',
      second: 'numeric',
      hour12: true
    };

    // Format the date to NYC time
    const nycTime = new Intl.DateTimeFormat('en-US', options).format(date);

    return nycTime;
  };

  const getFirstTime = (subway: SubwayData | null) => {
    if (!subway || !subway.entity || subway.entity.length === 0) {
      return '-';
    }
    return toNYCTime(subway?.entity[0]?.trip_update?.stop_time_update[0]?.arrival?.time);
  }

  return (
      <main className={styles.main}>
        <Weather data={weather} />
        <Forecast data={weather} />
        <Subway />
        <SubwaySouth />
        {/* <pre className='text-xs bg-blue-500 text-white inline-block rounded-lg p-0 pl-1 pr-1'>
        <strong>Page Connection Status:</strong> {ReadyState[readyState]}
      </pre>
      <div className='mt-3'>
        <strong>Weather:</strong> 
        <div>{weather ? weather.current.weather[0].description : '' }</div>
        <div>{weather ? weather.current.temp : '' }</div>
      </div>
      <div className='mt-3'>
        <strong>Subway Line A:</strong> 
        <div>{getFirstTime(subwayA)}</div>
      </div>
      <div className='mt-3'>
        <strong>Subway Line B:</strong> 
        <div>{getFirstTime(subwayB)}</div>
      </div>
      <div className='mt-3'>
        <strong>Subway Line C:</strong> 
        <div>{getFirstTime(subwayC)}</div>
      </div> */}
      </main>
  );
};

export default Main;
