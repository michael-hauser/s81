import moment from 'moment-timezone';

// Define the structure for arrival and departure times
interface TimeUpdate {
    time: number;
}

// Define the structure for stop time updates
interface StopTimeUpdate {
    stop_id: string;
    arrival: TimeUpdate;
    departure: TimeUpdate;
    departure_occupancy_status: number;
}

// Define the structure for trip details
interface Trip {
    trip_id: string;
    route_id: string;
    start_time: string;
    start_date: string;
}

// Define the structure for trip updates
interface TripUpdate {
    trip: Trip;
    stop_time_update: StopTimeUpdate[];
}

// Define the structure for each entity
interface Entity {
    id: string;
    trip_update: TripUpdate;
}

// Define the structure for the entire response
export interface SubwayData {
    entity: Entity[];
}

export enum Direction {
    North = 'N',
    South = 'S'
}

// Define the structure for subway arrival times for the React Widgets
export interface SubwayArrival {
    tripId: string;
    line: string;
    arrivalMinutes: number;
    direction: Direction;
}

// Function to convert seconds to minutes
const convertSecondsToMinutes = (seconds: number): number => {
    return Math.ceil(seconds / 60); // Round up to the nearest minute
}

export const mapSubwayData = (data: SubwayData): SubwayArrival[] => {
    const currentTime = Math.floor(moment().tz('America/New_York').unix());

    if (!data.entity) return [];

    return data.entity.flatMap(entity => {
        const tripUpdate = entity.trip_update;
        if (!tripUpdate) return [];
        let line = tripUpdate.trip.route_id;
        if(line === 'D') line = 'B';
        if (line !== 'A' && line !== 'B' && line !== 'C') return [];

        return tripUpdate.stop_time_update.map(stopTimeUpdate => {
            const arrivalTime = stopTimeUpdate.arrival.time;
            const arrivalMinutes = convertSecondsToMinutes(arrivalTime - currentTime);

            return {
                tripId: tripUpdate.trip.trip_id,
                line,
                arrivalMinutes: Math.max(arrivalMinutes, 0), // Ensure non-negative values
                direction: stopTimeUpdate.stop_id.includes('N') ? Direction.North : Direction.South
            };
        });
    })
    //filter out 0 arrival times and arrival times greater than 30 minutes
    .filter((arrival) => arrival.arrivalMinutes > 0 && arrival.arrivalMinutes < 30)
    .sort((a, b) => a.arrivalMinutes - b.arrivalMinutes);
};