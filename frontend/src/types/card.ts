export interface EmployeeSummary {
  id: number;
  name: string;
  department: string;
}

export interface CompanyValueSummary {
  id: number;
  code: string;
  name: string;
  type: 'VALUE' | 'CREDO';
}

export interface CompanyValueDetail extends CompanyValueSummary {
  description: string;
}

export interface EmojiSummary {
  emojiCode: string;
  count: number;
  userIds: number[];
}

export interface CardResponse {
  id: number;
  sender: EmployeeSummary;
  recipients: EmployeeSummary[];
  reason: string;
  values: CompanyValueSummary[];
  createdAt: string;
  reactions: EmojiSummary[];
}

export interface CreateCardRequest {
  recipientIds: number[];
  valueIds: number[];
  reason: string;
}

export interface CardListResponse {
  items: CardResponse[];
  total: number;
}

export interface CardFilters {
  page?: number;
  pageSize?: number;
  from?: string;
  to?: string;
  q?: string;
  values?: number[];
  senderIds?: number[];
  recipientIds?: number[];
}
