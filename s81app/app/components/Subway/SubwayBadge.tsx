import React, { memo, useEffect } from 'react';
import styles from './SubwayBadge.module.scss';

interface SubwayBadgeProps {
    line: string;
}

const SubwayBadge: React.FC<SubwayBadgeProps> = ({ line = 'A' }) => {  
    return (
        <div className={styles.badge} style={{ backgroundColor: `var(--subway${line})` }}>
            {line}
        </div>
    )
}

export default memo(SubwayBadge);