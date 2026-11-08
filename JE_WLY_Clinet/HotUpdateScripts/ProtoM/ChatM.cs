using JEngine.Event;
using JEngine.Helper;
using JEngine.Singleton;

using pb;

using System;
using System.Collections.Generic;

using static pb.ChatHistoryResponse;

namespace HotUpdateScripts.ProtoM
{
    /// <summary>
    /// 聊天
    /// </summary>
    public class ChatM : SingletonTemplate<ChatM>
    {
        //聊天
        private int MaxCount = 50;
        private List<ChatInfo> HallChatList = new List<ChatInfo>();
        private List<ChatInfo> FightChatList = new List<ChatInfo>(3);
        private List<ChatInfo> AllChatList = new List<ChatInfo>();

        public void ClearAllChat()
        {
            HallChatList.Clear();
            AllChatList.Clear();
            FightChatList.Clear();
        }

        public void AddOneChat(ChatPush data)
        {
            var chat = new ChatInfo { ChatType = data.Type, Context = data.Context, Name = data.Name, Uid = data.Uid, Time = GameUtil.Instance.GetCurrentTimeSeconds() };
            AddHallChat(chat);
            AllHallChat(chat);
        }

        public void AddChatHistory(List<ChatChannelHistory> chatHistories)
        {
            foreach (var chatHistory in chatHistories)
            {
                //Debug.Log($"ChatM AddChatHistory: Channel = {chatHistory.Channel} - count = {chatHistory.History.Count} chats");
                foreach (var chat in chatHistory.History)
                {
                    var chatInfo = new ChatInfo
                    {
                        ChatType = (ChatType)chatHistory.Channel,
                        Context = chat.Msg,
                        Name = chat.Name,
                        Uid = chat.Uid,
                        Time = chat.Time
                    };
                    AddHallChat(chatInfo);
                    AllHallChat(chatInfo);
                }
            }

            HallChatList.Sort((x, y) => x.Time.CompareTo(y.Time));
            AllChatList.Sort((x, y) => x.Time.CompareTo(y.Time));
        }


        /// <summary>
        /// 添加大厅聊天信息。最多显示两条聊天记录。
        /// </summary>
        /// <param name="data"></param>
        private void AddHallChat(ChatInfo data)
        {
            if (data.ChatType == ChatType.SYSTEM || (data.ChatType == ChatType.WORLD && Convert.ToBoolean(GameUtil.Instance.GetPlayerPrefsInt(GameConst.ChatOption.WorldMSG.ToString(), 1)))
                || (data.ChatType == ChatType.GUILD && Convert.ToBoolean(GameUtil.Instance.GetPlayerPrefsInt(GameConst.ChatOption.GuildMSG.ToString(), 1))))
            {
                HallChatList.Add(data);
                if (HallChatList.Count > 1)
                {
                    HallChatList.RemoveAt(0);
                }
                SetFightChat(data);
            }

            if (data.ChatType == ChatType.ZONE)
                SetFightChat(data);
        }

        private void SetFightChat(ChatInfo data)
        {
            if (FightChatList.Count >= 3)
                FightChatList.RemoveAt(0);
            FightChatList.Add(data);
        }

        /// <summary>
        /// 战场聊天缩略
        /// </summary>
        /// <returns></returns>
        public List<ChatInfo> GetFightChats()
        {
            return FightChatList;
        }

        /// <summary>
        /// 根据大厅需要显示的聊天类型条件，重新加载对应类型的聊天记录。
        /// </summary>
        public void LoadHallChat()
        {
            HallChatList.Clear();
            for (int i = 0; i < AllChatList.Count; i++)
            {
                AddHallChat(AllChatList[i]);
            }

            JEvent.defaultEvent.Post(new GameEvent.Sync_ChatPush { isOnlyHall = true });
        }

        private void AllHallChat(ChatInfo data)
        {
            if (data.ChatType != ChatType.PROMPT)
            {
                AllChatList.Add(data);
                if (AllChatList.Count > MaxCount)
                {
                    AllChatList.RemoveAt(0);
                }
            }
        }

        public List<ChatInfo> GetChatsByType(ChatShowType chatShowType)
        {
            switch (chatShowType)
            {
                case ChatShowType.World:
                    return AllChatList.FindAll(x => x.ChatType == ChatType.WORLD || x.ChatType == ChatType.SYSTEM);
                case ChatShowType.All:
                    return AllChatList.FindAll(x => x.ChatType == ChatType.WORLD || x.ChatType == ChatType.SYSTEM || x.ChatType == ChatType.GUILD);
                case ChatShowType.GUILD:
                    return AllChatList.FindAll(x => x.ChatType == ChatType.GUILD);
                case ChatShowType.Zone:
                    return AllChatList.FindAll(x => x.ChatType == ChatType.ZONE);
                default:
                    return new List<ChatInfo>();
            }
        }

        /// <summary>
        /// 获取大厅需要显示的聊天记录。
        /// </summary>
        /// <returns></returns>
        public List<ChatInfo> GetLastTwoChatOfAllType()
        {
            return HallChatList;
        }

        public bool CheckIsEmoji(string showStr)
        {
            if (showStr.Length > 9)
                return false;
            if (!showStr.Contains("emoji_"))
                return false;

            showStr = showStr.Replace("emoji_", "");
            // 尝试将字符串转换为整数
            bool isSuccess = int.TryParse(showStr, out _);

            return isSuccess;
        }

        public (string, string) GetChatTypeImg(ChatType chatType)
        {
            string icon = "";

            switch (chatType)
            {
                case ChatType.INVALID_ChatType:
                    break;
                case ChatType.SYSTEM:
                    icon = "img_xitong";
                    break;
                case ChatType.WORLD:
                    icon = "img_shijie";
                    break;
                case ChatType.GUILD:
                    icon = "img_hanghui";
                    break;
                case ChatType.SCROLLING:
                    break;
                case ChatType.PROMPT:
                    break;
                case ChatType.ZONE:
                    icon = "icon_Chat_zq";
                    break;
                default:
                    break;
            }

            return (GameConst.ChatAtlas, icon);
        }

        public string GetChatByType(ChatType chatType, string des, string title)
        {

            string nameColor = "";
            string chatContentColor = "";
            switch (chatType)
            {
                case ChatType.SYSTEM:
                    nameColor = "3e63af";
                    chatContentColor = "4188ba";
                    break;
                case ChatType.WORLD:
                    nameColor = "9b4417";
                    chatContentColor = "836c40";
                    break;
                case ChatType.GUILD:
                    nameColor = "3b7b23";
                    chatContentColor = "599d3a";
                    break;
                case ChatType.ZONE:
                    nameColor = "843285";
                    chatContentColor = "905098";
                    break;
                default:
                    break;
            }

            if (CheckIsEmoji(des))
                return $"<size=23><color=#{nameColor}>{title}:</color></size><color=#ffffff00>\n</color>";
            else
                return $"<size=23><color=#{nameColor}>{title}:</color></size><size=19><color=#{chatContentColor}>{des}</color></size>";
        }

        public enum ChatShowType
        {
            None,   // 无效
            World,  // 世界
            All,    // 综合
            GUILD,  // 行会
            Zone,    // 战区
        }

        public struct ChatInfo
        {
            public ChatType ChatType;
            public string Context;
            public string Name;
            public long Uid;
            public long Time;
        }
    }
}

