import React, { memo } from 'react';
import { View, Text, StyleSheet, ScrollView } from 'react-native';
import WeatherIcon from './WeatherIcon';
import Widget from './Widget';
import { ThemedText } from './ThemedText';
import { Colors } from '@/constants/Colors';
import { AreaChart, Grid, XAxis } from 'react-native-svg-charts';
import * as shape from 'd3-shape';
import { Defs, LinearGradient, Stop } from 'react-native-svg';
import { useColorScheme } from '@/hooks/useColorScheme';

interface ForecastProps {
  data?: WeatherData;
}

const Forecast: React.FC<ForecastProps> = ({ data }) => {
  const colorScheme = useColorScheme();

  const forecast = data?.daily || [];

  const getDayString = (date: number, long: boolean = false): string => {
    const d = new Date(date * 1000);
    const shortDayNames = ['Sun', 'Mon', 'Tue', 'Wed', 'Thu', 'Fri', 'Sat'];
    const longDayNames = ['Sunday', 'Monday', 'Tuesday', 'Wednesday', 'Thursday', 'Friday', 'Saturday'];
    return long ? longDayNames[d.getDay()] : shortDayNames[d.getDay()];
  };

  const chartData = forecast.map(d => d.temp.day);

  const contentInset = { top: 20, bottom: 20 };

  const themeFontColor = Colors[colorScheme ?? 'light'].text;
  const themeFontColorMuted = Colors[colorScheme ?? 'light'].mutedText;

  return (
    <Widget style={{ height: 250 }}>
      <ThemedText type='bold'>Forecast</ThemedText>
      {data && (
        <ScrollView horizontal style={styles.scrollContainer} contentContainerStyle={styles.scrollContent}>
          <View style={styles.chartContainer}>
            <AreaChart
              style={styles.areaChart}
              data={chartData}
              contentInset={contentInset}
              curve={shape.curveNatural}
              svg={{
                fill: 'url(#chartGradient)',
                stroke: 'url(#chartGradient)',
                strokeWidth: 2,
              }}
            >
              <Defs key={'defs'}>
                <LinearGradient id={'chartGradient'} x1={'0'} y1={'0'} x2={'0'} y2={'1'}>
                  <Stop offset={'0%'} stopColor={Colors.light.sunshine} stopOpacity={0.7} />
                  <Stop offset={'19%'} stopColor={Colors.light.sunshine} stopOpacity={0} />
                </LinearGradient>
              </Defs>
            </AreaChart>
          </View>
          <View style={styles.legend}>
            {forecast.map((day, index) => (
              <View key={index} style={styles.forecastDay}>
                <ThemedText style={[styles.forecastDate]}>{getDayString(day.dt)}</ThemedText>
                <WeatherIcon size={32} weatherData={day} />
                <View style={styles.forecastTemp}>
                  <Text style={[styles.day, { color: themeFontColor }]}>{Math.round(day.temp.day)}°</Text>
                  <Text style={[styles.night, { color: themeFontColorMuted }]}>{Math.round(day.temp.night)}°</Text>
                </View>
              </View>
            ))}
          </View>
        </ScrollView>
      )}
    </Widget>
  );
};

const styles = StyleSheet.create({
  scrollContainer: {
    flexDirection: 'row',
    flex: 1,
  },
  scrollContent: {
    flexDirection: 'column',
    alignItems: 'flex-start',
  },
  chartContainer: {
    flexDirection: 'column',
    flex: 1,
    height: 200,
  },
  areaChart: {
    flex: 1,
    height: 100,
    width: 500, // Increase this width to allow scrolling
  },
  xAxis: {
    marginHorizontal: -10,
    width: 500, // Match the width of the chart for consistency
  },
  legend: {
    flexDirection: 'row',
    marginTop: 16,
    alignItems: 'center',
  },
  forecastDay: {
    flexDirection: 'column',
    alignItems: 'center',
    marginHorizontal: 8,
  },
  forecastDate: {
    marginBottom: 4
  },
  forecastTemp: {
    flexDirection: 'row',
    width: '100%',
    marginTop: 4,
    gap: 4,
  },
  day: {
    fontSize: 14,
  },
  night: {
    fontSize: 14,
  },
});

export default memo(Forecast);
