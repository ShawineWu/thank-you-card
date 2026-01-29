#!/usr/bin/env python3
"""
Mock Data Generator for Thank You Card System
Generates SQL insert statements for hackathon demo
"""

import random
import uuid
from datetime import datetime, timedelta

# User data - Contains both Hydra ID and Microsoft ID
# hydraId: from hydra-test.castlery.com (used by web-dashboard login)
# msId: Microsoft Azure AD Object ID (used by Teams app)
# For users without hydraId yet, we use msId as the primary id
USERS = [
    {"id": "06715c3c-10b1-7000-9147-ee26ed2fc4c3", "msId": "d3dc438c-1de2-49dc-ae8b-8dffb7fcb2cc", "name": "Youshan Li (SZX)", "email": "youshan.li@castlery.com"},
    {"id": "02129f90-2224-486b-859f-b60c2da6b826", "msId": "02129f90-2224-486b-859f-b60c2da6b826", "name": "Zehui Lin (SZX)", "email": "zehui.lin@castlery.com"},
    {"id": "cd7e9fb0-d844-441e-aac4-b8c21d36ef7c", "msId": "cd7e9fb0-d844-441e-aac4-b8c21d36ef7c", "name": "Christopher Li (SZX)", "email": "christopher.li@castlery.com"},
    {"id": "3ef2d96d-35b0-42ee-8ace-e4e7880f12c0", "msId": "3ef2d96d-35b0-42ee-8ace-e4e7880f12c0", "name": "Shawine Wu (SZX)", "email": "shawine.wu@castlery.com"},
    {"id": "f1f7264d-5a6c-4fef-82e9-a122f9a332d0", "msId": "f1f7264d-5a6c-4fef-82e9-a122f9a332d0", "name": "Alex Zu (SZX)", "email": "alex.zu@castlery.com"},
    {"id": "dfaa0ed8-2473-403b-8fdb-7fbbb38dfd06", "msId": "dfaa0ed8-2473-403b-8fdb-7fbbb38dfd06", "name": "Bingye Li (SZX)", "email": "bingye.li@castlery.com"},
]

# Company values (from migration)
COMPANY_VALUES = [
    "MAKE_IMPACT", "STRIVE_EXCELLENCE", "STAND_TOGETHER", "BE_OPEN_MINDED", "STAY_GROUNDED",
    "BIAS_FOR_ACTION", "CUSTOMER_CENTRIC", "THINK_STRATEGICALLY", "DEEP_DIVE", "INVENT_SIMPLIFY",
    "EARN_TRUST", "TAKE_OWNERSHIP", "CHALLENGE_DISAGREE_COMMIT", "LEARN_BE_CURIOUS", "DO_MORE_WITH_LESS"
]

# Recognition reason templates (Chinese)
REASON_TEMPLATES = [
    "非常感谢{name}在{project}项目中的出色表现！{detail}你的专业能力和敬业精神让整个团队都受益匪浅。{extra}",
    "感谢{name}在过去一段时间里的辛勤付出！{detail}你的工作态度和专业素养值得我们每个人学习。{extra}",
    "{name}，谢谢你在{project}中的帮助！{detail}正是因为有你这样优秀的同事，我们的团队才能不断进步。{extra}",
    "特别感谢{name}在{project}期间的支持和协助！{detail}你的贡献对项目的成功至关重要。{extra}",
    "向{name}表达诚挚的感谢！{detail}你在{project}中展现出的专业能力令人印象深刻。{extra}",
    "感谢{name}一直以来的支持！{detail}在{project}项目中，你的付出让我们看到了团队协作的力量。{extra}",
    "{name}，非常感谢你在{project}中的贡献！{detail}你的努力和专注让整个项目得以顺利完成。{extra}",
    "衷心感谢{name}！{detail}在{project}项目中，你展现出了卓越的领导力和执行力。{extra}",
]

PROJECTS = [
    "Q1产品迭代", "系统重构", "性能优化", "新功能开发", "用户体验改进",
    "数据分析", "技术架构升级", "客户支持", "团队建设", "流程优化",
    "代码审查", "文档完善", "测试覆盖", "部署自动化", "监控告警"
]

DETAILS = [
    "你在技术方案设计上的深入思考，帮助我们避免了很多潜在问题。",
    "你主动承担了很多额外的工作，确保项目按时交付。",
    "你的代码质量一直保持很高的水准，为团队树立了榜样。",
    "你在遇到困难时从不退缩，总是积极寻找解决方案。",
    "你的沟通能力很强，能够清晰地表达技术方案和想法。",
    "你对细节的关注让我们的产品质量得到了显著提升。",
    "你总是愿意分享知识，帮助团队成员共同成长。",
    "你的创新思维为项目带来了很多新的可能性。",
    "你在压力下依然保持冷静，高效地完成了任务。",
    "你的团队协作精神让整个项目进展得更加顺利。",
]

EXTRAS = [
    "期待未来能有更多合作的机会！",
    "希望我们能继续保持这样的合作默契！",
    "你是团队中不可或缺的一员！",
    "感谢你为团队带来的正能量！",
    "你的努力不会被忽视，继续加油！",
    "很高兴能和你一起工作！",
    "你的专业精神值得我们每个人学习！",
    "期待你在未来取得更大的成就！",
    "感谢你一直以来的支持和帮助！",
    "你是我们团队的骄傲！",
    "希望我们能一起创造更多的价值！",
    "你的贡献对团队意义重大！",
]

# Milestone thresholds
MILESTONE_THRESHOLDS = [1, 5, 10, 25, 50, 100]

MILESTONE_TITLES = {
    "CardsSent": {
        1: "初次感谢",
        5: "感恩之心",
        10: "传递温暖",
        25: "感谢达人",
        50: "感恩大使",
        100: "感谢之星"
    },
    "CardsReceived": {
        1: "首次被认可",
        5: "团队之光",
        10: "优秀伙伴",
        25: "卓越贡献者",
        50: "团队明星",
        100: "传奇人物"
    }
}

MILESTONE_DESCRIPTIONS = {
    "CardsSent": {
        1: "发送了第一张感谢卡，开启了感恩之旅",
        5: "已发送5张感谢卡，持续传递正能量",
        10: "已发送10张感谢卡，成为团队中的温暖使者",
        25: "已发送25张感谢卡，是团队中的感谢达人",
        50: "已发送50张感谢卡，成为感恩文化的推动者",
        100: "已发送100张感谢卡，是团队中的感谢之星"
    },
    "CardsReceived": {
        1: "收到了第一张感谢卡，被同事认可",
        5: "已收到5张感谢卡，是团队中的闪光点",
        10: "已收到10张感谢卡，是大家信赖的伙伴",
        25: "已收到25张感谢卡，是团队中的卓越贡献者",
        50: "已收到50张感谢卡，是团队中的明星成员",
        100: "已收到100张感谢卡，是团队中的传奇人物"
    }
}

def generate_reason(recipient_name: str) -> str:
    """Generate a recognition reason between 100-500 characters"""
    template = random.choice(REASON_TEMPLATES)
    project = random.choice(PROJECTS)
    detail = random.choice(DETAILS)
    extra = random.choice(EXTRAS)
    
    reason = template.format(
        name=recipient_name.split(" (")[0],  # Use first name only
        project=project,
        detail=detail,
        extra=extra
    )
    
    # Ensure length is between 100-500
    while len(reason) < 100:
        reason += " " + random.choice(EXTRAS)
    
    if len(reason) > 500:
        reason = reason[:497] + "..."
    
    return reason

def random_timestamp(start_date: datetime, end_date: datetime) -> datetime:
    """Generate random timestamp between start and end date"""
    delta = end_date - start_date
    random_seconds = random.randint(0, int(delta.total_seconds()))
    return start_date + timedelta(seconds=random_seconds)

def escape_sql_string(s: str) -> str:
    """Escape single quotes for SQL"""
    return s.replace("'", "''")

def generate_mock_data():
    """Generate all mock data"""
    # Date range: 2025-01-26 to 2025-01-30
    start_date = datetime(2025, 1, 26, 0, 0, 0)
    end_date = datetime(2025, 1, 30, 23, 59, 59)
    
    # Track cards sent and received per user
    cards_sent = {user["id"]: 0 for user in USERS}
    cards_received = {user["id"]: 0 for user in USERS}
    
    # Target counts (random 10-100)
    target_sent = {user["id"]: random.randint(10, 100) for user in USERS}
    target_received = {user["id"]: random.randint(10, 100) for user in USERS}
    
    # Store all cards data
    cards_data = []
    card_recipients_data = []
    card_values_data = []
    
    # Generate cards to meet targets
    all_cards_needed = []
    
    # First, create cards based on sender targets
    for sender in USERS:
        sender_id = sender["id"]
        num_cards = target_sent[sender_id]
        
        for _ in range(num_cards):
            # Pick random recipient (not self)
            possible_recipients = [u for u in USERS if u["id"] != sender_id]
            recipient = random.choice(possible_recipients)
            
            all_cards_needed.append({
                "sender": sender,
                "recipient": recipient,
                "timestamp": random_timestamp(start_date, end_date)
            })
    
    # Sort by timestamp
    all_cards_needed.sort(key=lambda x: x["timestamp"])
    
    # Generate SQL data
    for card_info in all_cards_needed:
        card_id = str(uuid.uuid4())
        sender = card_info["sender"]
        recipient = card_info["recipient"]
        timestamp = card_info["timestamp"]
        
        reason = generate_reason(recipient["name"])
        
        # Card
        cards_data.append({
            "id": card_id,
            "sender_id": sender["id"],
            "sender_name": sender["name"],
            "recognition_reason": reason,
            "created_at": timestamp
        })
        
        # Card recipient
        card_recipients_data.append({
            "id": str(uuid.uuid4()),
            "card_id": card_id,
            "recipient_id": recipient["id"],
            "recipient_name": recipient["name"],
            "created_at": timestamp
        })
        
        # Card values (1-3 random values)
        num_values = random.randint(1, 3)
        selected_values = random.sample(COMPANY_VALUES, num_values)
        for value_code in selected_values:
            card_values_data.append({
                "id": str(uuid.uuid4()),
                "card_id": card_id,
                "value_code": value_code,
                "created_at": timestamp
            })
        
        # Update counts
        cards_sent[sender["id"]] += 1
        cards_received[recipient["id"]] += 1
    
    # Generate milestones based on actual counts
    milestones_data = []
    
    for user in USERS:
        user_id = user["id"]
        
        # Sent milestones
        sent_count = cards_sent[user_id]
        for threshold in MILESTONE_THRESHOLDS:
            if sent_count >= threshold:
                milestones_data.append({
                    "id": str(uuid.uuid4()),
                    "employee_id": user_id,
                    "milestone_type": "CardsSent",
                    "threshold": threshold,
                    "title": MILESTONE_TITLES["CardsSent"][threshold],
                    "description": MILESTONE_DESCRIPTIONS["CardsSent"][threshold],
                    "achieved_at": random_timestamp(start_date, end_date)
                })
        
        # Received milestones
        received_count = cards_received[user_id]
        for threshold in MILESTONE_THRESHOLDS:
            if received_count >= threshold:
                milestones_data.append({
                    "id": str(uuid.uuid4()),
                    "employee_id": user_id,
                    "milestone_type": "CardsReceived",
                    "threshold": threshold,
                    "title": MILESTONE_TITLES["CardsReceived"][threshold],
                    "description": MILESTONE_DESCRIPTIONS["CardsReceived"][threshold],
                    "achieved_at": random_timestamp(start_date, end_date)
                })
    
    return cards_data, card_recipients_data, card_values_data, milestones_data, cards_sent, cards_received

def generate_sql():
    """Generate SQL insert statements"""
    cards, recipients, values, milestones, sent_counts, received_counts = generate_mock_data()
    
    sql_lines = []
    sql_lines.append("-- Mock Data for Thank You Card System")
    sql_lines.append("-- Generated for Hackathon Demo")
    sql_lines.append(f"-- Date Range: 2025-01-26 to 2025-01-30")
    sql_lines.append(f"-- Total Cards: {len(cards)}")
    sql_lines.append("")
    
    # Print statistics
    sql_lines.append("-- Statistics:")
    for user in USERS:
        sql_lines.append(f"-- {user['name']}: Sent={sent_counts[user['id']]}, Received={received_counts[user['id']]}")
    sql_lines.append("")
    
    # Clear existing data
    sql_lines.append("-- Clear existing data")
    sql_lines.append("DELETE FROM card_values;")
    sql_lines.append("DELETE FROM card_recipients;")
    sql_lines.append("DELETE FROM employee_milestones;")
    sql_lines.append("DELETE FROM cards;")
    sql_lines.append("")
    
    # Insert cards
    sql_lines.append("-- Insert cards")
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
    sql_lines.append("")
    
    # Insert milestones
    sql_lines.append("-- Insert employee milestones")
    for milestone in milestones:
        ts = milestone["achieved_at"].strftime("%Y-%m-%d %H:%M:%S")
        title = escape_sql_string(milestone["title"])
        desc = escape_sql_string(milestone["description"])
        sql_lines.append(
            f"INSERT INTO employee_milestones (id, employee_id, milestone_type, threshold, achieved_at, title, description) VALUES "
            f"('{milestone['id']}', '{milestone['employee_id']}', '{milestone['milestone_type']}', {milestone['threshold']}, '{ts}', '{title}', '{desc}');"
        )
    
    return "\n".join(sql_lines)

if __name__ == "__main__":
    sql = generate_sql()
    
    # Write to file
    with open("scripts/mock_data.sql", "w", encoding="utf-8") as f:
        f.write(sql)
    
    print("Mock SQL generated successfully!")
    print("Output file: scripts/mock_data.sql")
