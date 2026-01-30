#!/usr/bin/env python3
"""
Weighted Mock Data Generator for Thank You Card System
Generates additional SQL insert statements to highlight specific company values
"""

import random
import uuid
from datetime import datetime, timedelta

# User data
USERS = [
    {"id": "06715c3c-10b1-7000-9147-ee26ed2fc4c3", "msId": "d3dc438c-1de2-49dc-ae8b-8dffb7fcb2cc", "name": "Youshan Li (SZX)", "email": "youshan.li@castlery.com"},
    {"id": "02129f90-2224-486b-859f-b60c2da6b826", "msId": "02129f90-2224-486b-859f-b60c2da6b826", "name": "Zehui Lin (SZX)", "email": "zehui.lin@castlery.com"},
    {"id": "cd7e9fb0-d844-441e-aac4-b8c21d36ef7c", "msId": "cd7e9fb0-d844-441e-aac4-b8c21d36ef7c", "name": "Christopher Li (SZX)", "email": "christopher.li@castlery.com"},
    {"id": "3ef2d96d-35b0-42ee-8ace-e4e7880f12c0", "msId": "3ef2d96d-35b0-42ee-8ace-e4e7880f12c0", "name": "Shawine Wu (SZX)", "email": "shawine.wu@castlery.com"},
    {"id": "f1f7264d-5a6c-4fef-82e9-a122f9a332d0", "msId": "f1f7264d-5a6c-4fef-82e9-a122f9a332d0", "name": "Alex Zu (SZX)", "email": "alex.zu@castlery.com"},
    {"id": "dfaa0ed8-2473-403b-8fdb-7fbbb38dfd06", "msId": "dfaa0ed8-2473-403b-8fdb-7fbbb38dfd06", "name": "Bingye Li (SZX)", "email": "bingye.li@castlery.com"},
]

# Highlighted values - these will be emphasized
HIGHLIGHTED_VALUES = [
    "BIAS_FOR_ACTION",
    "CUSTOMER_CENTRIC", 
    "STRIVE_EXCELLENCE"
]

# Recognition reason templates focused on highlighted values
REASON_TEMPLATES_BY_VALUE = {
    "BIAS_FOR_ACTION": [
        "非常感谢{name}在{project}项目中展现出的行动力！你总是第一时间响应需求，快速推动项目进展。",
        "感谢{name}的果断行动！在{project}中，你没有犹豫，立即采取行动解决问题，展现了出色的执行力。",
        "{name}，谢谢你在{project}中的快速响应！你的行动力让整个团队都受到了鼓舞。",
        "特别感谢{name}在{project}期间的高效执行！你用实际行动证明了什么是真正的行动派。",
    ],
    "CUSTOMER_CENTRIC": [
        "非常感谢{name}在{project}项目中始终以客户为中心！你总是站在客户角度思考问题。",
        "感谢{name}对客户需求的深刻理解！在{project}中，你的客户导向思维帮助我们做出了正确的决策。",
        "{name}，谢谢你在{project}中展现的客户至上精神！你总是把客户体验放在第一位。",
        "特别感谢{name}在{project}期间对客户反馈的重视！你的客户思维让产品更加贴近用户需求。",
    ],
    "STRIVE_EXCELLENCE": [
        "非常感谢{name}在{project}项目中追求卓越的精神！你对质量的坚持让整个团队都受益匪浅。",
        "感谢{name}对卓越的不懈追求！在{project}中，你的高标准推动了整个项目的品质提升。",
        "{name}，谢谢你在{project}中展现的卓越精神！你永远不满足于现状，总是追求更好。",
        "特别感谢{name}在{project}期间的精益求精！你的卓越标准成为了团队的榜样。",
    ],
}

PROJECTS = [
    "Q1产品迭代", "系统重构", "性能优化", "新功能开发", "用户体验改进",
    "数据分析", "技术架构升级", "客户支持", "团队建设", "流程优化",
]

EXTRAS = [
    "期待未来能有更多合作的机会！",
    "希望我们能继续保持这样的合作默契！",
    "你是团队中不可或缺的一员！",
    "感谢你为团队带来的正能量！",
    "你的努力不会被忽视，继续加油！",
    "很高兴能和你一起工作！",
]

def generate_reason(recipient_name: str, value_code: str) -> str:
    """Generate a recognition reason for specific value"""
    templates = REASON_TEMPLATES_BY_VALUE.get(value_code, REASON_TEMPLATES_BY_VALUE["BIAS_FOR_ACTION"])
    template = random.choice(templates)
    project = random.choice(PROJECTS)
    extra = random.choice(EXTRAS)
    
    reason = template.format(
        name=recipient_name.split(" (")[0],
        project=project
    ) + " " + extra
    
    return reason

def random_timestamp(start_date: datetime, end_date: datetime) -> datetime:
    """Generate random timestamp between start and end date"""
    delta = end_date - start_date
    random_seconds = random.randint(0, int(delta.total_seconds()))
    return start_date + timedelta(seconds=random_seconds)

def escape_sql_string(s: str) -> str:
    """Escape single quotes for SQL"""
    return s.replace("'", "''")

def generate_weighted_data():
    """Generate additional data to highlight specific values"""
    # Date range: same as original
    start_date = datetime(2025, 1, 26, 0, 0, 0)
    end_date = datetime(2025, 1, 30, 23, 59, 59)
    
    cards_data = []
    card_recipients_data = []
    card_values_data = []
    
    # Generate 150 additional cards focused on highlighted values
    # Distribution: BIAS_FOR_ACTION: 60, CUSTOMER_CENTRIC: 50, STRIVE_EXCELLENCE: 40
    value_counts = {
        "BIAS_FOR_ACTION": 60,
        "CUSTOMER_CENTRIC": 50,
        "STRIVE_EXCELLENCE": 40
    }
    
    for value_code, count in value_counts.items():
        for _ in range(count):
            # Pick random sender and recipient
            sender = random.choice(USERS)
            possible_recipients = [u for u in USERS if u["id"] != sender["id"]]
            recipient = random.choice(possible_recipients)
            
            card_id = str(uuid.uuid4())
            timestamp = random_timestamp(start_date, end_date)
            reason = generate_reason(recipient["name"], value_code)
            
            cards_data.append({
                "id": card_id,
                "sender_id": sender["id"],
                "sender_name": sender["name"],
                "recognition_reason": reason,
                "created_at": timestamp
            })
            
            card_recipients_data.append({
                "id": str(uuid.uuid4()),
                "card_id": card_id,
                "recipient_id": recipient["id"],
                "recipient_name": recipient["name"],
                "created_at": timestamp
            })
            
            # Primary highlighted value
            card_values_data.append({
                "id": str(uuid.uuid4()),
                "card_id": card_id,
                "value_code": value_code,
                "created_at": timestamp
            })
            
            # 30% chance to add a second highlighted value
            if random.random() < 0.3:
                other_values = [v for v in HIGHLIGHTED_VALUES if v != value_code]
                second_value = random.choice(other_values)
                card_values_data.append({
                    "id": str(uuid.uuid4()),
                    "card_id": card_id,
                    "value_code": second_value,
                    "created_at": timestamp
                })
    
    return cards_data, card_recipients_data, card_values_data

def generate_sql():
    """Generate SQL insert statements"""
    cards, recipients, values = generate_weighted_data()
    
    sql_lines = []
    sql_lines.append("-- Additional Weighted Mock Data for Thank You Card System")
    sql_lines.append("-- Generated to highlight: BIAS_FOR_ACTION, CUSTOMER_CENTRIC, STRIVE_EXCELLENCE")
    sql_lines.append(f"-- Date Range: 2025-01-26 to 2025-01-30")
    sql_lines.append(f"-- Additional Cards: {len(cards)}")
    sql_lines.append("")
    
    # Count values
    value_counts = {}
    for v in values:
        code = v["value_code"]
        value_counts[code] = value_counts.get(code, 0) + 1
    
    sql_lines.append("-- Value Distribution:")
    for code, count in sorted(value_counts.items(), key=lambda x: -x[1]):
        sql_lines.append(f"-- {code}: {count}")
    sql_lines.append("")
    
    # Insert cards
    sql_lines.append("-- Insert additional cards")
    for card in cards:
        reason = escape_sql_string(card["recognition_reason"])
        sender_name = escape_sql_string(card["sender_name"])
        ts = card["created_at"].strftime("%Y-%m-%d %H:%M:%S")
        sql_lines.append(
            f"INSERT INTO cards (id, sender_id, sender_name, recognition_reason, created_at, updated_at) VALUES "
            f"('{card['id']}', '{card['sender_id']}', '{sender_name}', '{reason}', '{ts}', '{ts}');"
        )
    sql_lines.append("")
    
    # Insert card recipients
    sql_lines.append("-- Insert card recipients")
    for recipient in recipients:
        ts = recipient["created_at"].strftime("%Y-%m-%d %H:%M:%S")
        name = escape_sql_string(recipient["recipient_name"])
        sql_lines.append(
            f"INSERT INTO card_recipients (id, card_id, recipient_id, recipient_name, created_at) VALUES "
            f"('{recipient['id']}', '{recipient['card_id']}', '{recipient['recipient_id']}', '{name}', '{ts}');"
        )
    sql_lines.append("")
    
    # Insert card values
    sql_lines.append("-- Insert card values")
    for value in values:
        ts = value["created_at"].strftime("%Y-%m-%d %H:%M:%S")
        sql_lines.append(
            f"INSERT INTO card_values (id, card_id, value_id, created_at) VALUES "
            f"('{value['id']}', '{value['card_id']}', (SELECT id FROM company_values WHERE code = '{value['value_code']}'), '{ts}');"
        )
    
    return "\n".join(sql_lines)

if __name__ == "__main__":
    sql = generate_sql()
    
    # Write to file
    with open("scripts/weighted_mock_data.sql", "w", encoding="utf-8") as f:
        f.write(sql)
    
    print("Weighted mock SQL generated successfully!")
    print("Output file: scripts/weighted_mock_data.sql")
    print("\nThis will ADD data to existing database, not replace it.")
