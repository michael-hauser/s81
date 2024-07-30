import React, { memo } from 'react';
import sharedStyles from '../Widget.module.scss';
import styles from './Weather.module.scss';
import WeatherIcon from '../WeatherIcon';

interface WeatherProps {
    data?: WeatherData
}

const Weather: React.FC<WeatherProps> = ({ data }) => {

    let temp = "";
    if (data && data.current) {
        temp = data.current.temp + "°F";
    }

    return (
        <div className={sharedStyles.widget}>
            <h2>Weather</h2>
            <div className={`${sharedStyles.widgetContent} ${styles.weatherContent}`}>
                <div className={styles.temp}>{temp}</div>
                <WeatherIcon size={150}/>
                </div>
        </div>
    )
}

export default memo(Weather);