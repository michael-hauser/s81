/**
 * Below are the colors that are used in the app. The colors are defined in the light and dark mode.
 * There are many other ways to style your app. For example, [Nativewind](https://www.nativewind.dev/), [Tamagui](https://tamagui.dev/), [unistyles](https://reactnativeunistyles.vercel.app), etc.
 */

const subwayA = '#103BA0';
const subwayB = '#EE6D35';
const subwayC = '#103BA0';
const sunshine = '#FFBA0A';

export const Colors = {
  light: {
    text: '#000000',
    mutedText: '#939393',
    background: '#F3F2EE',
    foreground: '#FFFFFF',
    hover: 'rgba(0, 0, 0, 0.05)',
    subwayA,
    subwayB,
    subwayC,
    sunshine,
  },
  dark: {
    text: '#FFFFFF',
    mutedText: '#939393',
    background: '#000000',
    foreground: '#212121',
    hover: 'rgba(255, 255, 255, 0.1)',
    subwayA,
    subwayB,
    subwayC,
    sunshine,
  },
};


const tintColorLight = '#0a7ea4';
const tintColorDark = '#fff';

export const OldColors = {
  light: {
    text: '#11181C',
    background: '#fff',
    tint: tintColorLight,
    icon: '#687076',
    tabIconDefault: '#687076',
    tabIconSelected: tintColorLight,
  },
  dark: {
    text: '#ECEDEE',
    background: '#151718',
    tint: tintColorDark,
    icon: '#9BA1A6',
    tabIconDefault: '#9BA1A6',
    tabIconSelected: tintColorDark,
  },
};