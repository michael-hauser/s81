import React, { memo } from 'react';
import sharedStyles from '../Widget.module.scss';
import styles from './Weather.module.scss';
import WeatherIcon from '../WeatherIcon';

interface WeatherProps {
    data?: WeatherData
}

const Weather: React.FC<WeatherProps> = ({ data }) => {
    const getTemp = (temp: number | undefined) => {
        if (!temp) return "";
        return Math.round(temp) + "°F";
    }

    return (
        <div className={sharedStyles.widget}>
            <h2>Weather</h2>
            <div className={`${sharedStyles.widgetContent} ${styles.weatherContent}`}>
                <div className={styles.temp}>{getTemp(data?.current?.temp)}</div>
                <WeatherIcon size={130} weatherData={data?.current}/>
            </div>
        </div>
    )
}

export default memo(Weather);