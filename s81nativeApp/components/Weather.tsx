import React, { memo } from 'react';
import { View, Text, StyleSheet } from 'react-native';
import WeatherIcon from './WeatherIcon';
import { ThemedView } from './ThemedView';
import { ThemedText } from './ThemedText';
import { BorderRadius } from '@/constants/Spaces';

interface WeatherProps {
  data?: WeatherData;
}

const Weather: React.FC<WeatherProps> = ({ data }) => {
  const getTemp = (temp: number | undefined) => {
    if (!temp) return "55 °F";
    return Math.round(temp) + "°F";
  };

  return (
    <ThemedView type='foreground' style={styles.widget}>
      <ThemedText style={styles.title}>Weather</ThemedText>
      <View style={styles.weatherContent}>
        <ThemedText style={styles.temp}>{getTemp(data?.current?.temp)}</ThemedText>
        <WeatherIcon size={130} weatherData={data?.current} />
      </View>
    </ThemedView>
  );
};

const styles = StyleSheet.create({
  widget: {
    borderRadius: BorderRadius.lg,
    display: 'flex',
    gap: 16,
    paddingVertical: 24,
    paddingHorizontal: 32,
    minHeight: 290,
    overflow: 'hidden',
  },
  title: {
  },
  weatherContent: {
    display: 'flex',
    flex: 1,
  },
  temp: {
    fontSize: 48,
    lineHeight: 48,
    fontWeight: '200',
  },
});

export default memo(Weather);
