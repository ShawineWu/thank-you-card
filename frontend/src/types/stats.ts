export interface ValueCount {
  valueId: number;
  code: string;
  name: string;
  type: 'VALUE' | 'CREDO';
  count: number;
}

export interface PersonalStats {
  totalSent: number;
  totalReceived: number;
  topSentValues: ValueCount[];
  topReceivedValues: ValueCount[];
}

export interface TopEmployee {
  employeeId: number;
  name: string;
  department: string;
  count: number;
}
