import { Colors } from '@/constants/Colors';
import { useColorScheme } from '@/hooks/useColorScheme';
import React, { memo } from 'react';
import { View, Text, StyleSheet } from 'react-native';

interface SubwayBadgeProps {
    line: string;
}

const SubwayBadge: React.FC<SubwayBadgeProps> = ({ line = 'A' }) => {
    const colorScheme = useColorScheme();
    
    const backgroundColor = (Colors[colorScheme ?? 'light'] as any)[`subway${line}`];

    return (
        <View style={[styles.badge, { backgroundColor }]}>
            <Text style={styles.text}>{line}</Text>
        </View>
    );
};

const styles = StyleSheet.create({
    badge: {
        borderRadius: 100,
        height: 28,
        width: 28,
        display: 'flex',
        justifyContent: 'center',
        alignItems: 'center',
    },
    text: {
        fontWeight: 'bold',
        color: '#fff',
    },
});

export default memo(SubwayBadge);
