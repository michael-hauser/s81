import React, { useState, useEffect } from 'react';
import { StyleSheet, Text, View, ScrollView } from 'react-native';
import Forecast from '@/components/Forecast';
import Weather from '@/components/Weather';
import { ThemedText } from '@/components/ThemedText';
import { ThemedView } from '@/components/ThemedView';
import Time from '@/components/Time';
import { Direction, mapSubwayData, SubwayArrival } from '@/models/subwayData';
import useWebSocket from 'react-use-websocket';
import Subway from '@/components/Subway';

type WebSocketMessage = {
  key: 'weather-data' | 'subway-a' | 'subway-b' | 'subway-c';
  value: string;
};

export default function Home() {
  const WS_URL = 'ws://192.168.1.78:8081/ws';
  const { lastJsonMessage } = useWebSocket(WS_URL, {
    share: true,
    shouldReconnect: () => true,
  });

  const [weather, setWeather] = useState<WeatherData | undefined>(undefined);
  const [subwayA, setSubwayA] = useState<SubwayArrival[]>([]);
  const [subwayB, setSubwayB] = useState<SubwayArrival[]>([]);
  const [subwayC, setSubwayC] = useState<SubwayArrival[]>([]);

  useEffect(() => {
    console.log('lastJsonMessage:', lastJsonMessage);
    if (lastJsonMessage) {
      const message: WebSocketMessage = lastJsonMessage as WebSocketMessage;
      console.log('Received message:', message);

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

  return (
    <ScrollView contentContainerStyle={styles.scrollContainer}>
      <ThemedView style={styles.app}>
        <ThemedView style={styles.header}>
          <ThemedView>
            <ThemedText type='title'>
              <Text>81st Street-Museum of </Text>
              <Text>Natural History station</Text>
            </ThemedText>
            <ThemedView style={styles.time}>
              <Time />
            </ThemedView>
          </ThemedView>
        </ThemedView>

        <ThemedView style={styles.main}>
          <Weather data={weather} />
          <Forecast data={weather} />
          <Subway arrivals={subwayData} direction={Direction.North} />
          <Subway arrivals={subwayData} direction={Direction.South} />
        </ThemedView>
      </ThemedView>
    </ScrollView>
  );
}

const styles = StyleSheet.create({
  scrollContainer: {
    paddingHorizontal: 16,
    paddingVertical: 24,
  },
  app: {
    flexDirection: 'column',
    gap: 20,
    width: '100%',
  },
  header: {
    flexDirection: 'row',
    alignItems: 'flex-start',
    justifyContent: 'space-between',
    paddingHorizontal: 8,
    paddingTop: 8,
  },
  time: {
    marginTop: 16,
  },
  main: {
    width: '100%',
    maxWidth: '100%',
    gap: 16,
  },
});
