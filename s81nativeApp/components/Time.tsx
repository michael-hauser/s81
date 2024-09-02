import React from 'react';
import { useEffect, useState } from "react";
import { ThemedText } from './ThemedText';

export default function Time() {
    const [time, setTime] = useState<string>('');

    useEffect(() => {
      const updateTime = () => {
        const now = new Date();
        const nycTime = new Intl.DateTimeFormat('en-US', {
          timeZone: 'America/New_York',
          weekday: 'long',
          year: 'numeric',
          month: 'long',
          day: 'numeric',
          hour: 'numeric',
          minute: 'numeric',
          second: 'numeric',
          hour12: true,
        }).format(now);
  
        setTime(nycTime);
      };
  
      updateTime(); // Set initial time
      const intervalId = setInterval(updateTime, 1000); // Update every second
  
      return () => clearInterval(intervalId); // Cleanup interval on component unmount
    }, []);

    return (
        <ThemedText>
            {time}
        </ThemedText>
    );
};