import { ThemedText } from '@/components/ThemedText';
import { ThemedView } from '@/components/ThemedView';
import Time from '@/components/Time';
import Weather from '@/components/Weather';
import React from 'react';
import { StyleSheet, Text, View } from 'react-native';

export default function Home() {

    return (
        <ThemedView style={styles.app}>
            <ThemedView style={styles.header}>
                <ThemedView>
                    <ThemedText type='title'>
                        <Text>81st Street-Museum of </Text>
                        <Text>Natural History station</Text>
                    </ThemedText>
                    <ThemedView style={styles.time}>
                        <Time />
                    </ThemedView>
                </ThemedView>
            </ThemedView>
            
            <ThemedView>
                <Weather />
            </ThemedView>
        </ThemedView>
    )
}

const styles = StyleSheet.create({
    app: {
        display: 'flex',
        flexDirection: 'column',
        gap: 50,
        paddingHorizontal: 24,
        paddingVertical: 24,
        width: '100%',
        height: '100%',
    },
    header: {
        display: 'flex',
        flexDirection: 'row',
        alignItems: 'flex-start',
        justifyContent: 'space-between',
    },
    time: {
        marginTop: 16,
    },
    main: {
        width: '100%',
        maxWidth: '100%',
        display: 'flex',
        minWidth: 0,
        gap: 40,
    }
});