import React, { memo } from 'react';
import Image from 'next/image';
import cloud from '/public/cloud.svg';
import lightning from '/public/lightning.svg';
import rain from '/public/rain.svg';
import snow from '/public/snow.svg';
import sun from '/public/sun.svg';
import sunCloud from '/public/suncloud.svg';
import suncloudrain from '/public/suncloudrain.svg';

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

    return (
        <Image src={iconSrc} width={`${size}`} alt="weather icon" />
    )
}

export default memo(WeatherIcon);