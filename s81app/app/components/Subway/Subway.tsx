import React, { memo } from 'react';
import sharedStyles from '../Widget.module.scss';
import styles from './Subway.module.scss';
import SubwayBadge from './SubwayBadge';
import { Direction, SubwayArrival } from '@/app/models/subwayData';

interface SubwayProps {
    arrivals?: SubwayArrival[];
    direction: Direction;
}

const Subway: React.FC<SubwayProps> = ({ arrivals = [], direction }) => {

    const getPercentage = (arrivalMinutes: number): number => {
        const maxMinutes = 60; // Adjust this value based on your needs
        return Math.min((arrivalMinutes / maxMinutes) * 100, 100);
    }

    return (
        <div className={sharedStyles.widget}>
            <h2>{
                direction === Direction.North ? 'Northbound' : 'Southbound'
            }</h2>
            <div className={`${sharedStyles.widgetContent} ${styles.subwayContent}`}>
                {
                    arrivals
                        .filter(arrival => arrival.direction === direction)
                        .map((arrival, index) => (

                            <div key={index} className={styles.subwayArrival}>
                                <SubwayBadge line={arrival.line} />
                                <span className={styles.time}>{arrival.arrivalMinutes}m</span>
                                <div className={styles.animation} style={{ '--width': `${getPercentage(arrival.arrivalMinutes)}%` } as React.CSSProperties}></div>
                            </div>
                        ))
                }
            </div>
        </div>
    )
}

export default memo(Subway);