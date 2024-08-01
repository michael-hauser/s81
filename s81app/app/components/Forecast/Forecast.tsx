import React, { memo } from 'react';
import sharedStyles from '../Widget.module.scss';
import styles from './Forecast.module.scss';
import WeatherIcon from '../WeatherIcon';
import { Area, AreaChart, ResponsiveContainer, Tooltip } from 'recharts';

interface ForecastProps {
    data?: WeatherData
}

const Forecast: React.FC<ForecastProps> = ({ data }) => {

    const forecast = data?.daily || [];

    const getDayString = (date: number, long: boolean = false): string => {
        const d = new Date(date * 1000);
        const shortDayNames = ['Sun', 'Mon', 'Tue', 'Wed', 'Thu', 'Fri', 'Sat'];
        const longDayNames = ['Sunday', 'Monday', 'Tuesday', 'Wednesday', 'Thursday', 'Friday', 'Saturday'];    
        return long ? longDayNames[d.getDay()] : shortDayNames[d.getDay()];
    }

    const chartData = forecast.map(d => ({
        day: getDayString(d.dt, true),
        tempNight: d.temp.night + '°',
        valueTxt: d.temp.day + '°',
        value: d.temp.day,
    }));

    return (
        <div className={sharedStyles.widget}>
            <h2>Forecast</h2>
            {
                data && <div className={`${sharedStyles.widgetContent} ${styles.forecastContent}`}>
                    <div className={styles.chart}>
                        <svg style={{ position: "fixed", visibility: "hidden" }}>
                            <defs>
                                <linearGradient id="chartGradient" x1="0" x2="0" y1="0" y2="100%" gradientUnits="userSpaceOnUse">
                                    <stop offset="0%" stopColor="var(--sunshine)" stopOpacity={1} />
                                    <stop offset="70%" stopColor="var(--sunshine)" stopOpacity={0} />
                                </linearGradient>
                            </defs>
                        </svg>
                        <ResponsiveContainer width="100%" height="100%">
                            <AreaChart
                                width={600}
                                height={400}
                                data={chartData}
                            >
                                <Tooltip cursor={{ fill: 'rgba(0, 0, 0, 0.05)' }} content={({ active, payload, label }) => {
                                    return <div className={styles.tooltip}>
                                        <div className={styles.tooltipTitle}>{payload?.[0]?.payload.day}</div>
                                        <div className={styles.temp}>
                                            <span className={styles.day}>{payload?.[0]?.payload.valueTxt}</span>
                                            <span className={styles.night}>{payload?.[0]?.payload.tempNight}</span>
                                        </div>
                                    </div>;
                                }} />
                                <Area cursor={"pointer"} type="monotone" dataKey="value" stroke="var(--sunshine)" fill="url(#chartGradient)" />
                            </AreaChart>
                        </ResponsiveContainer>
                    </div>
                    <div className={styles.legend}>
                        {
                            forecast.map((day, index) => (
                                <div key={index} className={styles.forecastDay}>
                                    <div className={styles.forecastDate}>
                                        {getDayString(day.dt)}
                                    </div>
                                    <WeatherIcon size={32} />
                                    <div className={styles.forecastTemp}>
                                        <span className={styles.day}>{Math.round(day.temp.day)}°</span>
                                        <span className={styles.night}>{Math.round(day.temp.night)}°</span>
                                    </div>
                                </div>
                            ))
                        }
                    </div>
                </div>
            }
        </div>
    )
}

export default memo(Forecast);