import React from 'react';
import sharedStyles from '../Widget.module.scss';
import styles from './Subway.module.scss';
import SubwayBadge from './SubwayBadge';

const Subway: React.FC = () => {

    const [subwayArrivals, setSubwayArrivals] = React.useState<SubwayArrival[]>([
        { line: 'A', arrivalMinutes: 2 },
        { line: 'B', arrivalMinutes: 5 },
    ]);
    
    return (
        <div className={sharedStyles.widget}>
            <h2>Southbound</h2>
            <div className={`${sharedStyles.widgetContent} ${styles.subwayContent}`}>
                {subwayArrivals.map((arrival, index) => (
                    <div key={index} className={styles.subwayArrival}>
                        <SubwayBadge line={arrival.line} />
                        <span>{arrival.arrivalMinutes}m</span>
                    </div>
                ))}
            </div>
        </div>
    )
}

export default Subway;