import React, { memo } from 'react';
import Image from 'next/image';
import sunCloud from '/public/weather-sun-cloud.svg';

interface WeatherIconProps {
    weatherData?: CurrentWeather;
    size?: number;
}

const WeatherIcon: React.FC<WeatherIconProps> = ({ weatherData, size = 16 }) => {

    let iconSrc;

    if (!weatherData) {
        iconSrc = sunCloud;
    } else {
        switch (weatherData.weather[0].icon) {
            case '01d':
                iconSrc = sunCloud;
            default:
                iconSrc = sunCloud;
                break;
        }
    }

    return (
        <Image src={iconSrc} width={`${size}`} alt="weather icon" />
    )
}

export default memo(WeatherIcon);