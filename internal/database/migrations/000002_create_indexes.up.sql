CREATE INDEX idx_orders_vendor_id ON orders(vendor_id);
CREATE INDEX idx_delay_reports_order_id ON delay_reports(order_id);
CREATE INDEX idx_delay_reports_created_at ON delay_reports(created_at);
CREATE INDEX idx_delay_checks_agent_id ON delay_checks(agent_id);
CREATE INDEX idx_delay_checks_report_id ON delay_checks(report_id);