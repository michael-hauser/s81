// Home.tsx
import React from 'react';
import { StyleSheet, Text, View, ScrollView } from 'react-native';
import Forecast from '@/components/Forecast';
import Weather from '@/components/Weather';
import { ThemedText } from '@/components/ThemedText';
import { ThemedView } from '@/components/ThemedView';
import Time from '@/components/Time';
import { Direction } from '@/models/subwayData';
import Subway from '@/components/Subway';
import { useIndex } from '@/hooks/useIndex';

export default function Home() {
  const { weather, subwayData } = useIndex();

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
