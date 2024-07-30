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
interface SubwayData {
    entity: Entity[];
}

// Define the structure for subway arrival times for the React Widgets
interface SubwayArrival {
    line: string;
    arrivalMinutes: number;
}