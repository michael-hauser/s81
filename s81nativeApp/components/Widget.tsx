import React from 'react';
import { View, StyleSheet, ViewProps } from 'react-native';
import { ThemedView } from './ThemedView';
import { BorderRadius } from '@/constants/Spaces';

interface WidgetProps extends ViewProps {
  children: React.ReactNode;
  type?: 'background' | 'foreground';
}

const Widget: React.FC<WidgetProps> = ({ children, style, ...otherProps }) => {
  return (
    <ThemedView type='foreground' style={[styles.widget, style]} {...otherProps}>
      {children}
    </ThemedView>
  );
};

const styles = StyleSheet.create({
  widget: {
    borderRadius: BorderRadius.lg,
    paddingVertical: 16,
    paddingHorizontal: 24,
    minHeight: 160,
    overflow: 'visible'
  },
});

export default Widget;
