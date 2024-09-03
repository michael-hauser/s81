import React, { memo } from 'react';
import { View, Text, StyleSheet } from 'react-native';
import WeatherIcon from './WeatherIcon';
import Widget from './Widget';
import { ThemedText } from './ThemedText';

interface WeatherProps {
  data?: WeatherData;
}

const Weather: React.FC<WeatherProps> = ({ data }) => {
  const getTemp = (temp: number | undefined) => {
    if (!temp) return "55 °F";
    return Math.round(temp) + "°F";
  };

  return (
    <Widget>
      <ThemedText type='bold'>Weather</ThemedText>
      <View style={styles.weatherContent}>
        <ThemedText style={styles.temp}>{getTemp(data?.current?.temp)}</ThemedText>
        <WeatherIcon size={80} weatherData={data?.current} />
      </View>
    </Widget>
  );
};

const styles = StyleSheet.create({
  weatherContent: {
    flexDirection: 'row',
    alignItems: 'center',
    justifyContent: 'space-around',
    flex: 1,
  },
  temp: {
    fontSize: 50,
    lineHeight: 50,
    fontWeight: '200',
  },
});

export default memo(Weather);
