import React, { memo } from 'react';
import { View, Text, StyleSheet } from 'react-native';
import SubwayBadge from './SubwayBadge';
import { Direction, MAX_DISPLAY_MINUTES, SubwayArrival } from '@/models/subwayData';
import { Colors } from '@/constants/Colors';
import { ThemedView } from './ThemedView';
import Widget from './Widget';
import { ThemedText } from './ThemedText';

interface SubwayProps {
    arrivals?: SubwayArrival[];
    direction: Direction;
}

const Subway: React.FC<SubwayProps> = ({ arrivals = [], direction }) => {

    const getPercentage = (arrivalMinutes: number): number => {
        const maxMinutes = MAX_DISPLAY_MINUTES;
        return Math.abs(Math.min((arrivalMinutes / maxMinutes) * 100, 100));
    }

    return (
        <Widget>
            <ThemedText type='bold'>
                {direction === Direction.North ? 'Northbound' : 'Southbound'}
            </ThemedText>
            <View style={styles.subwayContent}>
                {arrivals
                    .filter(arrival => arrival.direction === direction)
                    .map((arrival, index) => (
                        <View key={index} style={styles.subwayArrival}>
                            <SubwayBadge line={arrival.line} />
                            <ThemedText style={styles.time}>{arrival.arrivalMinutes}m</ThemedText>
                            <View style={[styles.animation, { marginLeft: `${getPercentage(arrival.arrivalSeconds / 60)}%` }]}>
                                <View style={styles.animationCircle} />
                            </View>
                        </View>
                    ))}
            </View>
        </Widget>
    );
}

const styles = StyleSheet.create({
    subwayContent: {
        marginTop: 16,
        flexDirection: 'column',
        gap: 16,
    },
    subwayArrival: {
        flexDirection: 'row',
        alignItems: 'center',
        justifyContent: 'space-between',
        gap: 8,
        display: 'flex',
        position: 'relative',
        width: 150,
    },
    time: {
        width: 40,
        fontSize: 14,
    },
    animation: {
        flex: 1,
        marginLeft: 10,
        height: 16,
        borderRadius: 3,
        opacity: 0.5,
    },
    animationCircle: {
        position: 'absolute',
        top: 0,
        left: 0,
        height: 16,
        width: 16,
        borderRadius: 100,
        backgroundColor: Colors.light.mutedText,
    },
});

export default memo(Subway);
