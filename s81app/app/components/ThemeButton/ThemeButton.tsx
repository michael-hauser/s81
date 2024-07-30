"use client";

import React, { useEffect } from 'react';
import Image from 'next/image';
import styles from './ThemeButton.module.scss';
import Dark from '/public/dark.svg';
import Light from '/public/light.svg';

const ThemeButton: React.FC = () => {

    const [theme, setTheme] = React.useState('light');

    const changeTheme = (theme: string) => {
        setTheme(theme);
        window.localStorage.setItem('theme', theme);
        document.documentElement.classList.remove('light', 'dark');
        document.documentElement.classList.add(theme);
    }

    useEffect(() => {
        const localTheme = window.localStorage.getItem('theme');
        localTheme && changeTheme(localTheme);
    }, []);

    const handleClick = () => {
        if(theme === 'light') changeTheme('dark');
        else changeTheme('light');
    }

    return (
        <button className={styles.button} onClick={handleClick}>
            { theme === 'light'
            ? <Image src={Dark} width="24" height="24" alt="dark"/>
            : <Image src={Light} width="24" height="24" alt="dark"/>
            }
        </button>
    )
}

export default ThemeButton;