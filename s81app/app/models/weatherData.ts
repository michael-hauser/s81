// Define the structure for weather conditions
interface WeatherCondition {
  id: number;
  main: string;
  description: string;
  icon: string;
}

// Define the structure for current weather data
interface CurrentWeather {
  dt: number;
  sunrise: number;
  sunset: number;
  temp: number;
  feels_like: number;
  pressure: number;
  humidity: number;
  dew_point: number;
  uvi: number;
  clouds: number;
  visibility: number;
  wind_speed: number;
  wind_deg: number;
  weather: WeatherCondition[];
}

// Define the structure for minutely data
interface MinutelyData {
  dt: number;
  precipitation: number;
}

// Define the structure for hourly weather data
interface HourlyWeather {
  dt: number;
  temp: number;
  feels_like: number;
  pressure: number;
  humidity: number;
  dew_point: number;
  uvi: number;
  clouds: number;
  visibility: number;
  wind_speed: number;
  wind_deg: number;
  wind_gust?: number; // Optional since it may not always be present
  weather: WeatherCondition[];
  pop: number;
}

// Define the structure for daily weather data
interface DailyWeather {
  dt: number;
  sunrise: number;
  sunset: number;
  moonrise: number;
  moonset: number;
  moon_phase: number;
  summary: string;
  temp: {
      day: number;
      min: number;
      max: number;
      night: number;
      eve: number;
      morn: number;
  };
  feels_like: {
      day: number;
      night: number;
      eve: number;
      morn: number;
  };
  pressure: number;
  humidity: number;
  dew_point: number;
  wind_speed: number;
  wind_deg: number;
  wind_gust: number;
  weather: WeatherCondition[];
  clouds: number;
  pop: number;
  rain?: number; // Optional since it may not always be present
  uvi: number;
}

// Define the structure for weather data response
interface WeatherData {
  lat: number;
  lon: number;
  timezone: string;
  timezone_offset: number;
  current: CurrentWeather;
  minutely: MinutelyData[];
  hourly: HourlyWeather[];
  daily: DailyWeather[];
}
