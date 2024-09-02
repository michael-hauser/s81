import React, { memo } from 'react';
import { Image, StyleSheet } from 'react-native';


const cloud = require('@/assets/images/cloud.svg');
const lightning = require('@/assets/images/lightning.svg');
const rain = require('@/assets/images/rain.svg');
const snow = require('@/assets/images/snow.svg');
const sun = require('@/assets/images/sun.svg');
const sunCloud = require('@/assets/images/suncloud.png');
const suncloudrain = require('@/assets/images/suncloudrain.svg');

interface WeatherIconProps {
  weatherData?: CurrentWeather | DailyWeather;
  size?: number;
}

const WeatherIcon: React.FC<WeatherIconProps> = ({ weatherData, size = 16 }) => {
  let iconSrc;

  if (!weatherData) {
    iconSrc = sunCloud;
  } else {
    switch (weatherData.weather[0].icon) {
      case '01d':
      case '01n':
        iconSrc = sun;
        break;
      case '02d':
      case '02n':
        iconSrc = sunCloud;
        break;
      case '03d':
      case '03n':
        iconSrc = cloud;
        break;
      case '04d':
      case '04n':
        iconSrc = cloud;
        break;
      case '09d':
      case '09n':
        iconSrc = rain;
        break;
      case '10d':
      case '10n':
        iconSrc = suncloudrain;
        break;
      case '11d':
      case '11n':
        iconSrc = lightning;
        break;
      case '13d':
      case '13n':
        iconSrc = snow;
        break;
      case '50d':
      case '50n':
        iconSrc = cloud;
        break;
      default:
        iconSrc = sunCloud;
        break;
    }
  }

  return <Image source={iconSrc} width={size} height={size} style={styles.icon} />;
};

const styles = StyleSheet.create({
  icon: {
    resizeMode: 'contain',
  },
});

export default memo(WeatherIcon);
