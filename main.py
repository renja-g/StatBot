import discord
import enum
from datetime import datetime
import asqlite
from dotenv import load_dotenv
import os

# This Bot will create a timeline for each member in the server


load_dotenv()
DISCORD_TOKEN = os.getenv('DISCORD_TOKEN')


intents = discord.Intents.default()
intents.guilds = True
intents.members = True
intents.presences = True
intents.voice_states = True

client = discord.Client(intents=intents)

DB_NAME = "data.db"


class Status(enum.Enum):
    ONLINE = 'online'
    ONLINE_MOBILE = 'online_mobile'
    IDLE = 'idle'
    IDLE_MOBILE = 'idle_mobile'
    DND = 'dnd'
    DND_MOBILE = 'dnd_mobile'
    OFFLINE = 'offline'


async def init_db():
    async with asqlite.connect(DB_NAME) as conn:
        async with conn.cursor() as c:
            await c.execute('''
                CREATE TABLE IF NOT EXISTS presence_update (
                    id INTEGER PRIMARY KEY AUTOINCREMENT,
                    user_id INTEGER,
                    status_after TEXT,
                    timestamp TEXT
                )
            ''')
            await c.execute('''
                CREATE TABLE IF NOT EXISTS activity_update (
                    id INTEGER PRIMARY KEY AUTOINCREMENT,
                    user_id INTEGER,
                    activity_after TEXT,
                    timestamp TEXT
                )
            ''')
            await c.execute('''
                CREATE TABLE IF NOT EXISTS voice_state_update (
                    id INTEGER PRIMARY KEY AUTOINCREMENT,
                    user_id INTEGER,
                    channel_id_after INTEGER,
                    channel_name_after TEXT,
                    deaf_after BOOLEAN,
                    mute_after BOOLEAN,
                    self_deaf_after BOOLEAN,
                    self_mute_after BOOLEAN,
                    self_stream_after BOOLEAN,
                    self_video_after BOOLEAN,
                    timestamp TEXT
                )
            ''')
            await conn.commit()


@client.event
async def on_ready():
    print(f'Logged in as {client.user}')
    await init_db()


def determine_status(member: discord.Member) -> Status:
    status_map = {
        (discord.Status.online, False): Status.ONLINE,
        (discord.Status.online, True): Status.ONLINE_MOBILE,
        (discord.Status.idle, False): Status.IDLE,
        (discord.Status.idle, True): Status.IDLE_MOBILE,
        (discord.Status.dnd, False): Status.DND,
        (discord.Status.dnd, True): Status.DND_MOBILE,
        (discord.Status.offline, False): Status.OFFLINE
    }

    return status_map.get((member.status, member.is_on_mobile()))


async def log_presence_update(user_id: int, status: Status):
    async with asqlite.connect(DB_NAME) as conn:
        async with conn.cursor() as c:
            await c.execute('''
                INSERT INTO presence_update (user_id, status_after, timestamp)
                VALUES (?, ?, ?)
            ''', (user_id, status.value, datetime.now().isoformat()))
            await conn.commit()


async def log_activity_update(user_id: int, activity_after: str):
    async with asqlite.connect(DB_NAME) as conn:
        async with conn.cursor() as c:
            await c.execute('''
                INSERT INTO activity_update (user_id, activity_after, timestamp)
                VALUES (?, ?, ?)
            ''', (user_id, activity_after, datetime.now().isoformat()))
            await conn.commit()


async def log_voice_state_update(user_id: int, after: discord.VoiceState):
    async with asqlite.connect(DB_NAME) as conn:
        async with conn.cursor() as c:
            await c.execute('''
                INSERT INTO voice_state_update (user_id, channel_id_after, channel_name_after, deaf_after, mute_after, self_deaf_after, self_mute_after, self_stream_after, self_video_after, timestamp)
                VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
            ''', (user_id, after.channel.id if after.channel else None, after.channel.name if after.channel else None, after.deaf, after.mute, after.self_deaf, after.self_mute, after.self_stream, after.self_video, datetime.now().isoformat()))
            await conn.commit()


@client.event
async def on_presence_update(before: discord.Member, after: discord.Member):
    # Status change
    if after.status != before.status or after.is_on_mobile() != before.is_on_mobile():
        print(f'{after.name}:\nStatus: {before.status} -> {after.status}\nMobile: {before.is_on_mobile()} -> {after.is_on_mobile()}\n')
        status = determine_status(after)
        await log_presence_update(after.id, status)

    # Activity change (for now only primary activity and only the name)
    if before.activity != after.activity:
        before_activity_name = before.activity.name if before.activity else None
        after_activity_name = after.activity.name if after.activity else None
        
        if before_activity_name != after_activity_name:
            print(f'{after.name}:\nActivity: {before_activity_name} -> {after_activity_name}\n')
            await log_activity_update(after.id, after_activity_name)


@client.event
async def on_voice_state_update(member: discord.Member, before: discord.VoiceState, after: discord.VoiceState):
    if before != after:
        # check if something changed
        if (
            before.channel != after.channel
            or before.deaf != after.deaf
            or before.mute != after.mute
            or before.self_deaf != after.self_deaf
            or before.self_mute != after.self_mute
            or before.self_stream != after.self_stream
            or before.self_video != after.self_video
        ):
            print(
                f"{member.name}:\nChannel: {before.channel} -> {after.channel}\nDeaf: {before.deaf} -> {after.deaf}\nMute: {before.mute} -> {after.mute}\nSelf Deaf: {before.self_deaf} -> {after.self_deaf}\nSelf Mute: {before.self_mute} -> {after.self_mute}\nSelf Stream: {before.self_stream} -> {after.self_stream}\nSelf Video: {before.self_video} -> {after.self_video}\n"
            )
            await log_voice_state_update(member.id, after)



client.run(DISCORD_TOKEN)
