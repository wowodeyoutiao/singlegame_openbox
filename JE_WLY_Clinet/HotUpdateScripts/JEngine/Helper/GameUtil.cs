using BestHTTP.SecureProtocol.Org.BouncyCastle.Asn1.Ocsp;

using DG.Tweening;

using HotUpdateScripts;
using HotUpdateScripts.ProtoM;

using JEngine.LitJson;
using JEngine.Singleton;

using MetaFramework.UI;

using pb;

using SuperScrollView;

using System;
using System.Collections.Generic;
using System.Text;

using Unity.Jobs;

using UnityEngine;
using UnityEngine.UI;

using static CFG_XMLConfig.XML_Config;

using Object = UnityEngine.Object;

namespace JEngine.Helper
{
    public class GameUtil : SingletonTemplate<GameUtil>
    {
        // Unix 纪元的 DateTime 对象 (1970 年 1 月 1 日)
        private DateTime epoch = new DateTime(1970, 1, 1, 0, 0, 0, DateTimeKind.Utc);
        public long GetCurrentTimeSeconds()
        {
            // 计算两个 DateTime 对象之间的时间差
            //TimeSpan timeSpan = DateTime.Now.ToUniversalTime() - epoch;
            //// 将 TimeSpan 对象转换为秒
            //long TotalSeconds = (long)timeSpan.TotalSeconds - GameM.Instance.SeverTimeOff;

            return (long)GameM.Instance.SeverTime;
        }

        /// 获取当前时间到指定至时间总毫秒秒数
        public double GetCurrentMSTimeSeconds()
        {
            // 计算两个 DateTime 对象之间的时间差
            TimeSpan timeSpan = DateTime.Now.ToUniversalTime() - epoch;
            return timeSpan.TotalMilliseconds - GameM.Instance.SeverTimeOff * 1000;
        }

        /// <summary>
        /// 时间戳转换为时间
        /// </summary>
        /// <param name="seconds">秒</param>
        /// <returns></returns>
        public DateTime ConvertFromUnixTimestampSeconds(long seconds)
        {
            return DateTimeOffset.FromUnixTimeSeconds(seconds).DateTime.ToLocalTime();
        }

        /// <summary>
        /// 时间戳转换为时间
        /// </summary>
        /// <param name="milliseconds">毫秒</param>
        /// <returns></returns>
        public DateTime ConvertFromUnixTimestampMilliseconds(long milliseconds)
        {
            return DateTimeOffset.FromUnixTimeMilliseconds(milliseconds).DateTime.ToLocalTime();
        }

        /// <summary>
        /// 将当前时间戳转换到当天0点的时间戳
        /// </summary>
        /// <param name="millisecondsTimestamp"></param>
        /// <returns></returns>
        public long ConvertToMidnightSeconds(long millisecondsTimestamp)
        {
            // 定义 Unix 时间的起始点（1970-01-01 00:00:00 UTC）
            DateTime unixEpoch = new DateTime(1970, 1, 1, 0, 0, 0, DateTimeKind.Utc);

            // 将毫秒时间戳转换为 DateTime
            DateTime dateTime = unixEpoch.AddMilliseconds(millisecondsTimestamp).ToLocalTime();

            // 将 DateTime 调整到零点
            DateTime midnightDateTime = dateTime.Date;

            // 将零点 DateTime 转换为秒级时间戳
            long secondsTimestampAtMidnight = (long)(midnightDateTime - unixEpoch.ToLocalTime()).TotalSeconds;

            return secondsTimestampAtMidnight;
        }

        public string GetTimeToYMD(long timestampInSeconds)
        {
            DateTime dateTimeUtc = DateTimeOffset.FromUnixTimeSeconds(timestampInSeconds).LocalDateTime;
            // 转换为"年-月-日 时：分"格式的字符串
            return dateTimeUtc.ToString("yyyy-MM-dd HH:mm");
        }

        public DateTime GetDataTime(long timestampInSeconds)
        {
            // 转换为"年-月-日 时：分"格式的字符串
            return DateTimeOffset.FromUnixTimeSeconds(timestampInSeconds).LocalDateTime;
        }

        /// <summary>
        /// 获取当前时间到指定至时间总秒数
        /// </summary>
        /// <returns></returns>
        public long GetNowToTimeSecoud(int day, double houre, int min = 0, int seound = 0)
        {
            DateTime dateTimeLocal = DateTimeOffset.FromUnixTimeSeconds(GetCurrentTimeSeconds()).LocalDateTime;

            DateTime nextDayMidnight = dateTimeLocal.Date;
            if (day > 0)
                nextDayMidnight += TimeSpan.FromDays(day);

            if (houre > 0)
                nextDayMidnight += TimeSpan.FromHours(houre);

            if (min > 0)
                nextDayMidnight += TimeSpan.FromMinutes(min);

            if (seound > 0)
                nextDayMidnight += TimeSpan.FromSeconds(min);

            TimeSpan timeSpanToEndOfTheDay = nextDayMidnight - dateTimeLocal;

            return (timeSpanToEndOfTheDay.Ticks / TimeSpan.TicksPerSecond);
        }

        /// <summary>
        /// 当前时间到角色注册后指定至时间段的总秒数
        /// </summary>
        /// <returns></returns>
        public long GetRegistRoleTimeSecoud(int day, int hour = 0, int min = 0, int seound = 0)
        {
            DateTime registTime = DateTimeOffset.FromUnixTimeSeconds(RoleM.Instance.RegisterAt).LocalDateTime;
            DateTime dateTimeLocal = DateTimeOffset.FromUnixTimeSeconds(GetCurrentTimeSeconds()).LocalDateTime;
            DateTime nextDayMidnight = registTime.Date;

            if (day > 0)
                nextDayMidnight += TimeSpan.FromDays(day);

            if (hour > 0)
                nextDayMidnight += TimeSpan.FromHours(hour);

            if (min > 0)
                nextDayMidnight += TimeSpan.FromMinutes(min);

            if (seound > 0)
                nextDayMidnight += TimeSpan.FromSeconds(min);

            TimeSpan timeSpanToEndOfTheDay = nextDayMidnight - dateTimeLocal;
            long seconds = timeSpanToEndOfTheDay.Ticks / TimeSpan.TicksPerSecond;
            return Math.Max(seconds, 0);
        }

        /// <summary>
        /// 距离角色注册后的第几天还剩几天？
        /// 注册当天算第一天
        /// </summary>
        /// <param name="day">角色注册后的第几天</param>
        /// <returns></returns>
        public int DaysRegistRoleTimeToNow(int day)
        {
            DateTime registTime = DateTimeOffset.FromUnixTimeSeconds(RoleM.Instance.RegisterAt).LocalDateTime;
            DateTime dateTimeLocal = DateTimeOffset.FromUnixTimeSeconds(GetCurrentTimeSeconds()).LocalDateTime;
            DateTime regTemp = registTime.Date;
            regTemp += TimeSpan.FromDays(day - 1);
            DateTime nowTemp = dateTimeLocal.Date;
            int days = (regTemp - nowTemp).Days;
            return Math.Max(days, 0);
        }

        /// <summary>
        /// 角色注册的第几天了
        /// </summary>
        /// <returns></returns>
        public int GetDayOfRegistRoleTime()
        {
            DateTime registTime = DateTimeOffset.FromUnixTimeSeconds(RoleM.Instance.RegisterAt).LocalDateTime;
            DateTime dateTimeLocal = DateTimeOffset.FromUnixTimeSeconds(GetCurrentTimeSeconds()).LocalDateTime;
            return (dateTimeLocal - registTime).Days + 1;
        }

        /// <summary>
        /// 已到达角色注册后的第几天
        /// </summary>
        /// <param name="day">角色注册后的第几天</param>
        /// <returns></returns>
        public bool IsArrivedRegistRoleTime(int day)
        {
            return DaysRegistRoleTimeToNow(day) == 0;
        }

        /// <summary>
        /// 获取时间已过去多久
        /// </summary>
        /// <returns></returns>
        public string GetTimePassedFromStep(long endTime)
        {
            DateTime nowTime = DateTimeOffset.FromUnixTimeSeconds(GetCurrentTimeSeconds()).LocalDateTime;
            DateTime compareTime = DateTimeOffset.FromUnixTimeSeconds(endTime).LocalDateTime;
            TimeSpan timeSpan = nowTime - compareTime;

            if (timeSpan.Days > 0)
                return $"{timeSpan.Days}天前";
            else if (timeSpan.Hours > 0)
                return $"{timeSpan.Hours}小时前";
            else if (timeSpan.Minutes > 0)
                return $"{timeSpan.Minutes}分钟前";

            return $"{timeSpan.Seconds}秒前";
        }

        /// <summary>
        /// 计算当前时间到截至时间的显示 hh:mm:ss
        /// </summary>
        /// <returns></returns>
        public string GetTimeSToHMS(long endTime)
        {
            TimeSpan timeSpan = TimeSpan.FromSeconds(endTime);

            return $"{timeSpan.Hours:00}:{timeSpan.Minutes:00}:{timeSpan.Seconds:00}";
        }

        /// <summary>
        /// 计算当前时间到截至时间的显示 天 hh:mm:ss
        /// </summary>
        /// <param name="endTime"></param>
        /// <returns></returns>
        public string GetTimeSToDHOrHMS(long endTime)
        {
            TimeSpan timeSpan = TimeSpan.FromSeconds(endTime);

            if (timeSpan.Days > 0)
            {
                if (timeSpan.Hours > 0)
                    return $"{timeSpan.Days}天{timeSpan.Hours}小时";
                else
                    return $"{timeSpan.Days}天";
            }
            else if (timeSpan.Hours > 0)
                return $"{timeSpan.Hours:00}:{timeSpan.Minutes:00}:{timeSpan.Seconds:00}";
            else
                return $"{timeSpan.Minutes:00}:{timeSpan.Seconds:00}";
        }

        /// <summary>
        /// 计算当前时间到截至时间的显示 天 hh:mm:ss
        /// </summary>
        /// <param name="endTime"></param>
        /// <returns></returns>
        public string GetTimeSToDHMS(long endTime)
        {
            TimeSpan timeSpan = TimeSpan.FromSeconds(endTime);

            if (timeSpan.Days > 0)
                return $"{timeSpan.Days}天{timeSpan.Hours:00}:{timeSpan.Minutes:00}:{timeSpan.Seconds:00}";
            else
                return $"{timeSpan.Hours:00}:{timeSpan.Minutes:00}:{timeSpan.Seconds:00}";
        }

        /// <summary>
        /// 转换为 "天 时" 或 "时 分" 格式的字符串。
        /// 小于一天的时间只显示时和分，大于等于一天的时间显示天和时。
        /// </summary>
        /// <param name="timestampInSeconds"></param>
        /// <returns></returns>
        public string GetTimeToDHorHM(long timestampInSeconds)
        {
            TimeSpan timeSpan = TimeSpan.FromSeconds(timestampInSeconds);
            if (timeSpan.Days > 0)
                return $"{timeSpan.Days}天{timeSpan.Hours}时";
            else
                return $"{timeSpan.Hours}时{timeSpan.Minutes}分";
        }

        /// <summary>
        /// 转换为"天 时 分"格式的字符串
        /// </summary>
        /// <param name="timestampInSeconds"></param>
        /// <returns></returns>
        public string GetTimeToDHM(long timestampInSeconds)
        {
            TimeSpan timeSpan = TimeSpan.FromSeconds(timestampInSeconds);
            int minutes = timeSpan.Seconds == 0 ? timeSpan.Minutes : timeSpan.Minutes + 1;
            return $"{timeSpan.Days}天{timeSpan.Hours}时{minutes}分";
        }

        /// <summary>
        /// 转换为"天 h：m"格式的字符串
        /// </summary>
        /// <param name="timestampInSeconds"></param>
        /// <returns></returns>
        public string GetTimeToDHM2(long timestampInSeconds)
        {
            TimeSpan timeSpan = TimeSpan.FromSeconds(timestampInSeconds);
            // 转换为"年-月-日 时：分"格式的字符串
            return $"{timeSpan.Days}天{timeSpan.Hours}:{timeSpan.Minutes}";
        }

        /// <summary>
        /// 转换为"天 hh:mm"格式的字符串
        /// </summary>
        /// <param name="timestampInSeconds"></param>
        /// <returns></returns>
        public string GetTimeToDHM3(long timestampInSeconds)
        {
            TimeSpan timeSpan = TimeSpan.FromSeconds(timestampInSeconds);
            return $"{timeSpan.Days}天{timeSpan.Hours:00}:{timeSpan.Minutes:00}";
        }

        public string GetTimeToDHM4(long timestampInSeconds)
        {
            TimeSpan timeSpan = TimeSpan.FromSeconds(timestampInSeconds);
            if (timeSpan.Days > 0)
                return $"{timeSpan.Days}天{timeSpan.Hours:00}:{timeSpan.Minutes:00}";
            else
                return $"{timeSpan.Hours:00}:{timeSpan.Minutes:00}:{timeSpan.Seconds:00}";
        }

        /// <summary>
        /// 获取时间天数
        /// </summary>
        /// <param name="timestampInSeconds"></param>
        /// <returns></returns>
        public int GetTimeToD(long timestampInSeconds)
        {
            return TimeSpan.FromSeconds(timestampInSeconds).Days;
        }

        /// <summary>
        /// 转换为"年-月-日 时：分"格式的字符串
        /// </summary>
        /// <param name="timestampInSeconds"></param>
        /// <returns></returns>
        public TimeSpan GetTimeSpan(long timestampInSeconds)
        {
            TimeSpan timeSpan = TimeSpan.FromSeconds(timestampInSeconds);
            return timeSpan;
        }

        /// <summary>
        /// 转换为 "时 分 秒" ， "分 秒" 或 "秒" 格式的字符串。
        /// </summary>
        /// <param name="endTime"></param>
        /// <returns></returns>
        public string GetTimeToHMS(long endTime)
        {
            TimeSpan timeSpan = TimeSpan.FromSeconds(endTime);

            if (timeSpan.TotalHours >= 1)
                return $"{Math.Floor(timeSpan.TotalHours)}时{timeSpan.Minutes}分{timeSpan.Seconds}秒";
            else if (timeSpan.Minutes > 0)
                return $"{timeSpan.Minutes}分{timeSpan.Seconds}秒";
            else
                return $"{timeSpan.Seconds}秒";
        }

        //检测是否在开始结束时间
        public bool CheakCurTime(string strTime, string endTime)
        {
            DateTime nowTime = DateTimeOffset.FromUnixTimeSeconds(GetCurrentTimeSeconds()).LocalDateTime;

            // 解析开始时间和结束时间
            TimeSpan startParseTime = DateTime.Parse(strTime).TimeOfDay;
            TimeSpan endParseTime = DateTime.Parse(endTime).TimeOfDay;

            // 将当前时间的时间部分提取出来
            TimeSpan nowTimeOfDay = nowTime.TimeOfDay;

            // 检查当前时间是否在开始时间和结束时间之间
            return startParseTime <= nowTimeOfDay && nowTimeOfDay < endParseTime;
        }

        //检测是否在开始结束时间
        public bool CheakCurTimeDay(string strTime, string endTime)
        {
            DateTime nowTime = DateTimeOffset.FromUnixTimeSeconds(GetCurrentTimeSeconds()).LocalDateTime;

            // 创建当天的开始时间和结束时间
            DateTime startTimeOfDay = nowTime.Date.Add(DateTime.Parse(strTime).TimeOfDay);
            DateTime endTimeOfDay = nowTime.Date.Add(DateTime.Parse(endTime).TimeOfDay);

            // 检查当前时间是否在开始时间和结束时间之间
            return nowTime >= startTimeOfDay && nowTime <= endTimeOfDay;
        }

        /// <summary>
        /// 检查是否在指定周几
        /// </summary>
        /// <param name="dayOfWeek">周几？</param>
        /// <returns></returns>
        public bool CheakCurTimeWeek(int dayOfWeek)
        {
            DateTime nowTime = DateTimeOffset.FromUnixTimeSeconds(GetCurrentTimeSeconds()).LocalDateTime;
            int _dayOfWeek = dayOfWeek == 7 ? 0 : dayOfWeek;
            if (nowTime.DayOfWeek == (DayOfWeek)_dayOfWeek)
                return true;
            return false;
        }

        public int CheckNextOpenTime(string strTime, string endTime)
        {
            DateTime _nowTime = DateTimeOffset.FromUnixTimeSeconds(GetCurrentTimeSeconds()).LocalDateTime;
            DateTime _startTime = _nowTime.Date.Add(DateTime.Parse(strTime).TimeOfDay);
            DateTime _endTime = _nowTime.Date.Add(DateTime.Parse(endTime).TimeOfDay);

            if (_startTime <= _nowTime && _nowTime < _endTime)
                return 0;
            else
            {
                if (_nowTime < _startTime)
                {
                    TimeSpan timeSpan = _startTime - _nowTime;
                    return (int)Math.Ceiling(timeSpan.TotalSeconds);
                }
                else
                {
                    DateTime nextStartTime = _startTime.AddDays(1);
                    return (int)(nextStartTime - _nowTime).TotalSeconds;
                }
            }
        }

        public int CheckNextEndTime(string endTime)
        {
            DateTime _nowTime = DateTimeOffset.FromUnixTimeSeconds(GetCurrentTimeSeconds()).LocalDateTime;
            DateTime _endTime = _nowTime.Date.Add(DateTime.Parse(endTime).TimeOfDay);

            if (_nowTime < _endTime)
            {
                TimeSpan timeSpan = _endTime - _nowTime;
                return (int)Math.Ceiling(timeSpan.TotalSeconds);
            }
            else
            {
                DateTime nextEndTime = _endTime.Date.AddDays(1);
                return (int)Math.Ceiling((nextEndTime - _nowTime).TotalSeconds);
            }
        }

        public int CheckNextOpenTime(string strTime, string endTime, int dayOfWeek)
        {
            DateTime _nowTime = DateTimeOffset.FromUnixTimeSeconds(GetCurrentTimeSeconds()).LocalDateTime;
            DateTime _compareTime = new DateTime(_nowTime.Year, _nowTime.Month, _nowTime.Day);
            DateTime _startParseTime = DateTime.Parse(strTime);
            DateTime _endParseTime = DateTime.Parse(endTime);
            DateTime _startTime = _compareTime.AddHours(_startParseTime.Hour).AddMinutes(_startParseTime.Minute).AddSeconds(_startParseTime.Second);
            DateTime _endTime = _compareTime.AddHours(_endParseTime.Hour).AddMinutes(_endParseTime.Minute).AddSeconds(_endParseTime.Second);

            int _dayOfWeek = dayOfWeek == 7 ? 0 : dayOfWeek;
            if (_startTime <= _nowTime && _nowTime < _endTime && _nowTime.DayOfWeek == (DayOfWeek)_dayOfWeek)
            {
                return 0;
            }
            else
            {
                if (_nowTime < _startTime && _dayOfWeek == (int)_nowTime.DayOfWeek)
                {
                    TimeSpan timeSpan = _startTime - _nowTime;
                    return (int)Math.Ceiling(timeSpan.TotalSeconds);
                }
                else
                {
                    int addDays = GetTotalDaysToDayOfWeek(_nowTime, (DayOfWeek)_dayOfWeek);
                    DateTime nextStartTime = _startTime.Date.AddDays(addDays).AddHours(_startTime.Hour).AddMinutes(_startTime.Minute);
                    return (int)Math.Ceiling((nextStartTime - _nowTime).TotalSeconds);
                }
            }
        }

        /// <summary>
        /// 距离下一次这个时间有多少秒？
        /// </summary>
        /// <returns></returns>
        public int GetNextTimeSecond(string time)
        {
            DateTime nowTime = DateTimeOffset.FromUnixTimeSeconds(GetCurrentTimeSeconds()).LocalDateTime;
            return GetNextTimeSecond(time, (int)nowTime.DayOfWeek);
        }

        /// <summary>
        /// 距离下一次这个时间有多少秒？
        /// </summary>
        /// <returns></returns>
        public int GetNextTimeSecond(string time, int dayOfWeek)
        {
            DateTime nowTime = DateTimeOffset.FromUnixTimeSeconds(GetCurrentTimeSeconds()).LocalDateTime;
            DateTime targetTime = nowTime.Date.Add(DateTime.Parse(time).TimeOfDay);

            int _dayOfWeek = dayOfWeek == 7 ? 0 : dayOfWeek;
            if (nowTime < targetTime && _dayOfWeek == (int)nowTime.DayOfWeek)
            {
                TimeSpan timeSpan = targetTime - nowTime;
                return (int)Math.Ceiling(timeSpan.TotalSeconds);
            }
            else
            {
                int addDays = GetTotalDaysToDayOfWeek(nowTime, (DayOfWeek)_dayOfWeek);
                DateTime nextStartTime = targetTime.Date.AddDays(addDays);
                return (int)Math.Ceiling((nextStartTime - nowTime).TotalSeconds);
            }
        }

        public string GetWeekDayName(int dayOfWeek)
        {
            int _dayOfWeek = dayOfWeek == 7 ? 0 : dayOfWeek;
            return GetWeekDayName((DayOfWeek)_dayOfWeek);
        }

        public string GetWeekDayName(DayOfWeek dayOfWeek)
        {
            return dayOfWeek switch
            {
                DayOfWeek.Sunday => "周日",
                DayOfWeek.Monday => "周一",
                DayOfWeek.Tuesday => "周二",
                DayOfWeek.Wednesday => "周三",
                DayOfWeek.Thursday => "周四",
                DayOfWeek.Friday => "周五",
                DayOfWeek.Saturday => "周六",
                _ => "",
            };
        }

        /// <summary>
        /// 距离周几还有多少天。
        /// 比如今天是周5，目标是周6，那么返回1天。如果目标是周1，那么返回3天。
        /// </summary>
        /// <param name="dateTime">要计算的日期</param>
        /// <param name="targetDay">目标时间</param>
        /// <returns></returns>
        public int GetTotalDaysToDayOfWeek(DateTime dateTime, DayOfWeek targetDay)
        {
            // 计算从当前日期到目标星期的偏移量
            int daysUntilTarget = ((int)targetDay - (int)dateTime.DayOfWeek + 7) % 7;

            if (daysUntilTarget == 0)
            {
                daysUntilTarget = 7;
            }

            return daysUntilTarget;
        }

        /// <summary>
        /// 数值统一化显示
        /// </summary>
        /// <param name="number"></param>
        /// <returns></returns>
        public string FormatNumWithUnit(double number)
        {
            string unit = string.Empty;
            if (number >= 1000000000000)
            {
                double num = Math.Floor(number / 1000000000000);
                unit = $"{num}兆";
            }
            else if (number >= 100000000000)
            {
                double num = Math.Floor(number / 100000000000);
                unit = $"{num}千亿";
            }
            else if (number >= 100000000) // 大于等于1亿
            {
                number /= 1000000;
                int num1 = (int)Math.Floor(number / 100);
                int num2 = (int)number % 100;
                if (num2 > 0)
                {
                    if (num2 < 10)
                        unit = $"{num1}.0{num2}";
                    else
                        unit = $"{num1}.{num2}";
                }
                else
                {
                    unit = $"{num1}";
                }
                unit += "亿";
            }
            else if (number >= 100000) // 大于等于10万
            {
                number /= 100;

                int num1 = (int)Math.Floor(number / 100);
                unit = $"{num1}";
                unit += "万";
            }
            else
            {
                unit = number.ToString(); // 其他情况，直接转换为字符串
            }

            return unit;
        }

        /// <summary>
        /// 战力数值统一化显示
        /// </summary>
        /// <param name="number"></param>
        /// <returns></returns>
        public string FormatPowerWithUnit(double number)
        {
            string unit = string.Empty;

            if (number >= 100000000) // 大于等于1亿(保留三位小数)
            {
                number /= 100000;
                int num1 = (int)Math.Floor(number / 1000);
                int num2 = (int)number % 1000;
                if (num2 > 0)
                {
                    if (num2 < 10)
                        unit = $"{num1}.00{num2}";
                    else if (num2 < 100)
                        unit = $"{num1}.0{num2}";
                    else
                        unit = $"{num1}.{num2}";
                }
                else
                {
                    unit = $"{num1}";
                }
                unit += "亿";
            }
            else if (number >= 10000000) // 大于等于1000万
            {
                int num1 = (int)Math.Floor(number / 10000);
                unit += $"{num1}万";
            }
            else
            {
                unit = number.ToString(); // 其他情况，直接转换为字符串
            }

            return unit;
        }


        /// <summary>
        /// 数值转为百分比保留一位小数显示
        /// </summary>
        /// <param name="avalue"></param>
        /// <returns></returns>
        public string GetPercentStr(float avalue)
        {
            return $"{(float)(Math.Truncate(avalue * 1000) / 10)}%"; // 格式化为保留一位小数
        }

        /// <summary>
        /// 检测道具是否充足
        /// </summary>
        /// <param name="avalue"></param>
        /// <returns></returns>
        public string GetCheakNumStr(int id, long num, string str = null)
        {
            if (str == null)
            {
                str = FormatNumWithUnit(num);
            }
            return ComCoinToCheck(id, num, null, false) ? str : "<color=#DA4E38>" + str + "</color>";
        }

        /// <summary>
        /// 设置人民币icon
        /// </summary>
        /// <param name="itemId"></param>
        /// <param name="go"></param>
        public void SetShopItemIcon(int needmoney, Image img = null)
        {
            if (img == null)
                return;

            if (needmoney == 100 || needmoney == 101)
            {
                ResouseHelper.Instance.LoadSpriteByAtlas(GameConst.CommonAtlas, (needmoney == 100) ? "icon_Activity_common_rmb" : "icon_Activity_common_bq", img);
            }
            else
                SetItemIcon(needmoney, img);
        }


        /// <summary>
        /// 根据道具ID 填充图片
        /// </summary>
        /// <param name="itemId"></param>
        /// <param name="go"></param>
        public void SetItemIcon(int itemId, GameObject go)
        {
            if (go.GetComponent<Image>() != null)
                SetItemIcon(itemId, go.GetComponent<Image>());
        }

        public void SetItemIcon(int itemId, Image img = null, Image bgImage = null, Image natureImg = null, int _Quality = -1)
        {
            var cfgItem = CFG_Item.GetById(itemId);
            if (img != null)
            {
                ResouseHelper.Instance.LoadSpriteByAtlas(cfgItem.Icon, cfgItem.Iconname, img);
            }
            if (bgImage != null)
            {
                SetItemBg((_Quality < 0 ? (int)cfgItem.Quality : _Quality), bgImage);
            }

            if (natureImg != null)
            {
                natureImg.enabled = cfgItem.NatureMark != 0;
                if (cfgItem.NatureMark != 0)
                {
                    ResouseHelper.Instance.LoadSpriteByAtlas(GameConst.CommonAtlas, "nature_icon_" + cfgItem.NatureMark, natureImg);
                }
            }
        }

        public void SetItemJobIcon(int cfgId, Image img)
        {
            SetJobIcon((int)CFG_Item.GetById(cfgId).Occupylimit, img);
        }

        public void SetJobIcon(int job, Image img)
        {
            if (img != null)
            {
                if (job == 0)
                {
                    img.enabled = false;
                    return;
                }
                img.enabled = true;
                ResouseHelper.Instance.LoadSpriteByAtlas(GameConst.CommonAtlas, "img_job" + job, img);
            }
        }

        public void SetTitleIcon(int itemId, Image img = null)
        {
            if (img != null)
            {
                if (itemId == 0)
                {
                    ResouseHelper.Instance.LoadSpriteByAtlas(GameConst.CommonAtlas, "bg_common_tm", img);
                    return;
                }
                var cfg = CFG_Title.GetById(itemId);
                ResouseHelper.Instance.LoadSpriteByAtlas(cfg.Path, cfg.Icon.ToString(), img);
            }
        }

        /// <summary>
        /// 根据怪物ID填充怪物头像
        /// </summary>
        /// <param name="mId"></param>
        /// <param name="img"></param>
        public void SetMonsterIcon(int mId, Image img = null)
        {
            var bossCFG = CFG_Monster.GetById(mId);
            if (img != null)
            {
                ResouseHelper.Instance.LoadSpriteByAtlas(bossCFG.Icon, bossCFG.Iconname, img);
            }
        }

        /// <summary>
        /// 玩家头像转载
        /// </summary>
        /// <param name="hUrl"></param>
        /// <param name="img"></param>
        public void SetHeadImg(string hUrl, Image img = null)
        {
            if (string.IsNullOrEmpty(hUrl))
            {
                ResouseHelper.Instance.LoadSpriteByAtlas(GameConst.CommonAtlas, "Com_img_touxiang_bg2", img);
            }
        }


        /// <summary>
        /// 根据品质填充图品
        /// </summary>
        /// <param name="Quilte"></param>
        /// <param name="image"></param>
        public void SetItemBg(int Quilte, Image image)
        {
            ResouseHelper.Instance.LoadSpriteByAtlas(GameConst.CommonAtlas, "bg_icon" + Quilte, image);
        }

        public void SetItemBg(int Quilte, GameObject image)
        {
            SetItemBg(Quilte, image.GetComponent<Image>());
        }

        /// <summary>
        /// 设置技能图片
        /// </summary>
        /// <param name="skillId"></param>
        /// <param name="image"></param>
        public void SetSkillBg(int skillId, Image image)
        {
            int idx = 1;
            if (skillId != 0 && (CFG_Skill.GetById(skillId).Skilltype == (int)CFG_Skill.SkillType.Super || CFG_Skill.GetById(skillId).Skilltype == (int)CFG_Skill.SkillType.SuperCtr))
            {
                idx = 2;
            }
            ResouseHelper.Instance.LoadSpriteByAtlas(GameConst.CommonAtlas, "icon_caowei" + idx, image);
        }

        /// <summary>
        /// 设置伙伴品质底框（圆形底框）
        /// </summary>
        /// <param name="quality">伙伴配置</param>
        /// <param name="image">image</param>
        public void SetPartnetQualityFrame(int quality, Image image)
        {
            string resName = quality switch
            {
                1 => "img_baise",
                2 => "img_lvse",
                3 => "img_lanse",
                4 => "img_zise",
                5 => "img_fense",
                6 => "img_chengse",
                7 => "img_hongse",
                8 => "img_yellow",
                _ => "img_weijiesuo2",
            };

            ResouseHelper.Instance.LoadSpriteByAtlas(GameConst.CommonAtlas, resName, image);
        }

        public string GetColorString(string str, string colorHex)
        {
            return $"<color=#{colorHex}>{str}</color>";
        }

        public string GetColorString(string str, Color color)
        {
            return $"<color=#{ColorUtility.ToHtmlStringRGB(color)}>{str}</color>";
        }

        public void SetColorText(Text text, string str, string colorHex)
        {
            text.text = GetColorString(str, colorHex);
        }

        public void SetColorText(Text text, string str, Color color)
        {
            text.text = GetColorString(str, color);
        }

        private string[] QualityColorDic =
        {
            "",//0
            "B7B5B1", //1 白色
            "4DAE1F", //2 绿色
            "64cdfe", //3 蓝色
            "e597ff", //4 紫色
            "ffa1a1", //5 粉色
            "fd9938", //6 橙色
            "ff5855", //7 红色
            "ebd940", //8 黄色
        };

        private string[] QualityOutlintColorDic =
        {
            "",//0
            "666666", //1 白色
            "0A5214", //2 绿色
            "36429c", //3 蓝色
            "551577", //4 紫色

            "911414", //5 粉色
            "783b06", //6 橙色
            "6f1c08", //7 红色
            "753411", //8 黄色
        };

        /// <summary>
        /// 返回道具名字（不同品质不同颜色）
        /// </summary>
        /// <param name="itemId"></param>
        public string GetItemName(int itemId, int Qulity = -1)
        {
            string equipStr = (CFG_Item.GetById(itemId).Type == 5 && Qulity > 0) ? CFG_XMLConfig.xmlConfig.EquipQualityDic[Qulity]._prefix : "";
            return $"<color=#{QualityColorDic[Qulity < 0 ? (int)CFG_Item.GetById(itemId).Quality : Qulity]}>{equipStr}{CFG_Item.GetById(itemId).Name}</color>";
        }

        public string GetStrByQuality(int qyality, string str)
        {
            return $"<color=#{QualityColorDic[qyality]}>{str}</color>";
        }

        public string GetQualityColorName(int quality)
        {
            return quality switch
            {
                1 => "白色",
                2 => "绿色",
                3 => "蓝色",
                4 => "紫色",
                5 => "粉色",
                6 => "橙色",
                7 => "红色",
                8 => "彩色",
                _ => "",
            };
        }

        private string[] FightQualityColorDic =
        {
            "", // 0
            "fff6f6",// 1
            "95ff40",// 2
            "5efcf5",// 3
            "f287ff",// 4
            "ffacac",// 5
            "ff942b",// 6
            "fe3c3c",// 7
            "fff712",// 8
        };

        public string GetFightStrByQuality(int qyality, string str)
        {
            return $"<color=#{FightQualityColorDic[qyality]}>{str}</color>";
        }

        public Color GetQualityOutlineColorByItemID(int itemId)
        {
            return GetQualityOutlineColor((int)CFG_Item.GetById(itemId).Quality);
        }

        public Color GetQualityOutlineColor(int quality)
        {
            return GetQualityOutlineColor(quality, new Color(0.6f, 0.5f, 0.4f));
        }

        public Color GetQualityOutlineColor(int quality, Color defaultColor)
        {
            if (quality < 1)
            {
                return defaultColor;
            }

            bool result = ColorUtility.TryParseHtmlString($"#{QualityOutlintColorDic[quality]}", out Color color);
            if (result)
                return color;
            else
                return defaultColor;
        }

        /// <summary>
        /// 游戏统一购买
        /// </summary>
        /// <param name="MainId"></param>
        /// <param name="SusseCallBack"></param>
        public void GameToBuy(long MainId, Action SusseCallBack, bool isFree = false, RefineSlot refineSlot = null)
        {
            GameToBuy(MainId, 1, SusseCallBack, isFree, refineSlot);
        }

        public void GameToBuy(long MainId, int mailNum, Action SusseCallBack, bool isFree = false, RefineSlot refineSlot = null)
        {
            var storeCfg = CFG_Store.GetById((int)MainId);

            var httpInfo = new HTTP_StoreOperateRequset()
            {
                goods_uuid = MainId,
                num = mailNum,
                refine_slot = refineSlot,
            };

            if (!isFree && SDKCtr.Instance.CurrentPlatform == SDKCtr.PlatformType.WebGL_WX)
            {
                //广告
                if (storeCfg.Limite.needmoney == 101)
                {
                    HttpC.ShuShuSendAdShow(new HTTP_WindowShushuAdShowRequest
                    {
                        ad_type = "激励广告",
                        ad_position = MainId.ToString(),
                    });

                    if (!GameM.Instance.SkipAD)
                    {
                        AudioManager.Instance.PauseBgMusic();
                        SDKCtr.Instance.ShowRewardedVideo(() =>
                        {
                            OnBuy();
                            AdClose(true);

                            AudioManager.Instance.ReSumeBgMusic();
                        },
                        () =>
                        {
                            UITool.ShowTips("广告观看取消");
                            AdClose(false);
                            AudioManager.Instance.ReSumeBgMusic();
                        });

                        void AdClose(bool completed)
                        {
                            HttpC.ShuShuSendAdClose(new HTTP_WindowShushuAdCloseRequest
                            {
                                ad_type = "激励广告",
                                ad_position = MainId.ToString(),
                                view_complete = completed,
                            });
                        }
                    }
                    else
                        OnBuy();
                }
                //充值
                else if (storeCfg.Limite.needmoney == 100)
                {
                    //UITool.ShowTips("功能暂未开放");
                    httpInfo.platform = "wc_" + PlayerPrefsExtend.GetString(SDKCtr.SaveKey.WX_Platform.ToString());

                    Debug.Log(" httpInfo.platform  " + httpInfo.platform);
                    var WXItem = new HTTP_StoreOperateWXItem
                    {
                        mode = "goods",
                        buyQuantity = mailNum,
                        env = 1,
                        currencyType = "CNY",
                        productId = MainId.ToString(),
                        goodsPrice = storeCfg.Limite.Num * mailNum * 100,
                    };
                    httpInfo.extra = JsonMapper.ToJson(WXItem);
                    if (httpInfo.platform == "wc_ios")
                        OnBuy();
                    else
                        HttpC.OnCheakSeeion((_r) =>
                        {
                            if (_r.code == 0)
                                OnBuy();
                            else
                                SDKCtr.Instance.OnWXWebLogin((code) =>
                                {
                                    GameM.Instance.WXCode = code;
                                    Debug.Log(" GameToBuy  " + GameM.Instance.WXCode);
                                    HttpC.OnGetOpenID(code, YesFun: (msg) =>
                                    {
                                        Debug.Log(" OnGetOpenID Buy 刷新  " + GameM.Instance.WXCode);
                                        OnBuy();
                                    });
                                }, null);
                        });
                }
                else
                {
                    OnBuy();
                }
            }
            else
            {
                OnBuy();
            }

            void OnBuy()
            {
                Action winCall = delegate
                {
                    if (storeCfg.Limite.needmoney == 100)
                    {
                        var goodsPrice = storeCfg.Limite.Num * mailNum * 100;
                        Debug.Log(" 购买价格：  " + goodsPrice);
                        SDKCtr.Instance.DNPurchase(goodsPrice);
                    }

                    StoreM.Instance.SetBuyNumAddByID(MainId, mailNum);
                    SusseCallBack?.Invoke();
                };
                long _BuyEventID = StoreM.Instance.SetBuyEvent(winCall);
                httpInfo.request_id = _BuyEventID;

                HttpC.OnGameToBuy(JsonMapper.ToJson(httpInfo), (_r) =>
                {
                    if (_r.code == 0)
                    {
                        if (string.IsNullOrEmpty(_r.signData) && string.IsNullOrEmpty(_r.prepay_internal_id))
                        {
                            AudioManager.Instance.PlaySound(AudioKey.Sound.UI_PurchaseGoods);
                            UITool.ShowTips("购买成功");
                            //winCall();
                            StoreM.Instance.OnDoneBuyEvent(_BuyEventID);
                        }
                        else
                        {
                            //购买支付的成功事件等待服务器推送再执行
                            if (!string.IsNullOrEmpty(_r.signData))
                                SDKCtr.Instance.OnWXWebPay(_r.paySig, _r.signData, _r.signature, null, delegate
                                {
                                    UITool.ShowTips("购买结束");
                                    StoreM.Instance.OnDoneBuyEvent(_BuyEventID, true);
                                });
                            else
                            {
                                string severStr = JsonMapper.ToJson(new HTTP_WXCustomerItem
                                {
                                    source_type = "order",
                                    out_trade_no = _r.prepay_internal_id
                                });
                                SDKCtr.Instance.OnWXCustomerService(severStr, null, delegate
                                {
                                    UITool.ShowTips("购买结束");
                                    StoreM.Instance.OnDoneBuyEvent(_BuyEventID, true);
                                });
                            }
                        }
                    }
                    else
                        UITool.ShowTips(_r.error);
                });
            }

        }

        public string SetAndCheckStr(string str, int id, int num)
        {
            bool isEnough = ComCoinToCheck(id, num, null, false);

            return !isEnough ? $"<color=#DA4E38>{str}</color>" : str;
        }

        //商城ID检测
        public bool ComCoinToCheck(int mailId, Action ClsoeFun, bool NeedJump = true, int mailNum = 1)
        {
            var cfgInfo = CFG_Store.GetById(mailId);

            return ComCoinToCheck(cfgInfo.Limite.needmoney, cfgInfo.Limite.Num * mailNum, ClsoeFun, NeedJump);
        }

        //公用货币不足跳转
        public bool ComCoinToCheck(int ItemId, long Num, Action ClsoeFun, bool NeedJump = true)
        {
            string strTip = $"当前{GetItemName(ItemId)}不足，是否前往获取?";

            var message = new UIMessageBox.CallBack
            {
                ShowStr = strTip,
            };

            bool isEnough = BagM.Instance.GetNumByCfgID(ItemId) >= Num;
            if (!NeedJump || isEnough)
                return isEnough;

            // 行会贡献 许愿点
            if (ItemId == (int)CFG_Item.ItemID.GuildContribution || ItemId == (int)CFG_Item.ItemID.GuildWish)
            {
                message.YesFun = () =>
                {
                    ClsoeFun?.Invoke();
                    UITool.Open<UIGuildDonation>();
                };
            }

            //点券
            else if (ItemId == (int)CFG_Item.ItemID.DianQuan)
            {
                message.YesFun = () =>
                {
                    ClsoeFun?.Invoke();
                    UITool.Open<UIShop>(null, CFG_Store.StoreType.RMB);
                };
            }

            //人民币(直接返回 拉支付)
            else if (ItemId == 100)
            {
                return true;
            }
            //直升宝石
            else if (ItemId >= 44 && ItemId <= 50)
            {
                message.YesFun = () =>
                {
                    ClsoeFun?.Invoke();
                    UITool.ShowTips("暂未开放");
                    //UITool.Open<UIHall>(null, UIHall.HallType.PVE);
                };
            }
            //战魂
            else if (ItemId == (int)CFG_Item.ItemID.ZhanHun)
            {
                if (GameRoomM.Instance._GameState != GameRoomM.GameState.No)
                {
                    UITool.ShowTips("已在战斗中");
                    return false;
                }
                else
                {
                    message.YesFun = () =>
                    {
                        if (CheakIsLock(CFG_XMLConfig.XML_Config.SwitchOpne.Play_PVP))
                            return;
                        ClsoeFun?.Invoke();
                        UITool.Open<UIFightChoice>(null, UIFightChoice.SelectType.Limit, UIFightChoice.ChanllengeType.PvP);
                    };
                }
            }
            //试炼凭证
            else if (ItemId == (int)CFG_Item.ItemID.SLPZ)
            {
                if (GameRoomM.Instance._GameState != GameRoomM.GameState.No)
                {
                    UITool.ShowTips("已在战斗中");
                    return false;
                }
                else
                {
                    message.YesFun = () =>
                    {
                        if (CheakIsLock(CFG_XMLConfig.XML_Config.SwitchOpne.Play_SlTower))
                            return;
                        ClsoeFun?.Invoke();
                        UITool.Open<UIFightChoice>(null, UIFightChoice.SelectType.Challange, UIFightChoice.ChanllengeType.SLTower);
                    };
                }
            }
            else if (CFG_Item.GetById(ItemId).Type == (decimal)CFG_Item.ItemType.PartnerPice)
                message.YesFun = () =>
                {
                    if (CheakIsLock(CFG_XMLConfig.XML_Config.SwitchOpne.Role_SkillAndCall))
                        return;

                    ClsoeFun?.Invoke();
                    HttpC.GetPrivilegeDetail(1, (_r) =>
                    {
                        UITool.Open<UIPartnerCall>(null, _r);
                    }, () => { UITool.Close<UILoading>(); });
                };

            // 书页
            else if (ItemId == (int)CFG_Item.ItemID.ShuYe)
                message.YesFun = () =>
                {
                    //if (CheakIsLock(CFG_XMLConfig.XML_Config.SwitchOpne.Shop_SL))
                    //    return;

                    ClsoeFun?.Invoke();
                    UITool.Open<UIShop>(null, CFG_Store.StoreType.Coupon);
                };
            // 技能卷轴
            else if (ItemId == (int)CFG_Item.ItemID.SkillJuanZhou)
                message.YesFun = () =>
                {
                    ClsoeFun?.Invoke();
                    UITool.Open<UISpecialOffer>(null, ActivityType.GiftPay);
                };
            else if (ItemId == (int)CFG_Item.ItemID.Strengthen || ItemId == (int)CFG_Item.ItemID.UpStar)
            {
                message.YesFun = () =>
                {
                    ClsoeFun?.Invoke();
                    UITool.Open<UIShop>(null, CFG_Store.StoreType.Coin);
                };
            }
            else
            {
                message.YesFun = () =>
                {
                    ClsoeFun?.Invoke();
                    UITool.Open<UIShop>(null, CFG_Store.StoreType.Coupon);
                };
            }
            UITool.Open<UIMessageBox>(null, message);

            return isEnough;
        }


        private float GetArrtValue(pb.EquipInfo eInfo, pb.Attribute attribute)
        {
            float _value = 0;
            if (eInfo != null && eInfo.Attributes != null)
            {
                foreach (var item in eInfo.Attributes)
                {
                    if (item.attr == attribute)
                    {
                        _value = item.Num;
                        break;
                    }
                }
            }
            return _value;
        }

        private float baseFM = 10000f; //万分比分母
        public int GetEquipPower(pb.ItemInfo ItemInfo, int Job, int _TeamIdx)
        {
            List<Stats> StatsList = _TeamIdx == 0 ? RoleM.Instance.GetRoleAllStats() : PartnerM.Instance.GetOnePartnerStats(PartnerM.Instance.GetOnePartnerSlotInfo(_TeamIdx).PartnerUid).Attributes;

            var skillList = _TeamIdx == 0 ? SkillM.Instance.GetSkillInfoList() : PartnerM.Instance.GetSkillInfoByPartnerId(PartnerM.Instance.GetOnePartnerSlotInfo(_TeamIdx).PartnerUid);


            float _skillCombat = 0;
            for (int i = 0; i < skillList.Count; i++)
            {
                var cfgSkillInfo = CFG_Skill.GetById(skillList[i].CfgId);
                _skillCombat += cfgSkillInfo.Powerfactor * (1 + skillList[i].Advanced);
            }

            int _basehp = 0; int _wuGong = 0; int _wuFang = 0; int _moGong = 0; int _moFang = 0; int _baoJi = 0;
            foreach (var item in StatsList)
            {
                if (item.attr == pb.Attribute.MAX_HEALTH)
                    _basehp = item.Num;
                else if (item.attr == pb.Attribute.PHYSICAL_ATTACK)
                    _wuGong = item.Num;
                else if (item.attr == pb.Attribute.MAGIC_ATTACK)
                    _moGong = item.Num;
                else if (item.attr == pb.Attribute.PHYSICAL_DEFENSE)
                    _wuFang = item.Num;
                else if (item.attr == pb.Attribute.MAGIC_DEFENSE)
                    _moFang = item.Num;
                else if (item.attr == pb.Attribute.CRITICAL_STRIKE)
                    _baoJi = item.Num;

            }
            var cfgInfo = CFG_Item.GetById(ItemInfo.CfgId);

            float equipValue = CFG_XMLConfig.xmlConfig.EquipQualityDic[(int)ItemInfo.Equip.Quality]._value;
            //生命
            float hp = (_basehp + (float)cfgInfo.HP * equipValue + GetArrtValue(ItemInfo.Equip, pb.Attribute.MAX_HEALTH)) * (1 + GetArrtValue(ItemInfo.Equip, pb.Attribute.HEALTH_PERCENTAGE) / baseFM);

            float wf = (_wuFang + (float)cfgInfo.Phydef * equipValue + GetArrtValue(ItemInfo.Equip, pb.Attribute.PHYSICAL_DEFENSE)) * (1 + GetArrtValue(ItemInfo.Equip, pb.Attribute.PHYSICAL_DEFENSE_PERCENTAGE) / baseFM);
            float mf = (_moFang + (float)cfgInfo.Magdef * equipValue + GetArrtValue(ItemInfo.Equip, pb.Attribute.MAGIC_DEFENSE)) * (1 + GetArrtValue(ItemInfo.Equip, pb.Attribute.MAGIC_DEFENSE_PERCENTAGE) / baseFM);

            float jsXS = GetArrtValue(ItemInfo.Equip, pb.Attribute.DMT_FINAL_REDUCE_PERCENTAGE) + GetArrtValue(ItemInfo.Equip, pb.Attribute.DMT_REDUCE_PERCENTAGE_PLAYER) + GetArrtValue(ItemInfo.Equip, pb.Attribute.DMT_REDUCE_PERCENTAGE_MONSTER);
            jsXS = jsXS > 7500 ? 7500 : jsXS;

            float js = GetArrtValue(ItemInfo.Equip, pb.Attribute.DMT_REDUCE_FIXED_PLAYER) + GetArrtValue(ItemInfo.Equip, pb.Attribute.DMT_REDUCE_FIXED_MONSTER);
            //最终减伤
            float zzjs = GetArrtValue(ItemInfo.Equip, pb.Attribute.DMT_FINAL_REDUCE_FIXED);

            //有效值
            float AllHP = (hp + CFG_XMLConfig.xmlConfig.BaseValue.Power.defvalue * (wf + mf) / 2f) / (1f - jsXS / baseFM) + js / 2f + zzjs;
            //主攻击值
            float zgjz = 0;
            if (Job % 10 == 1)
                zgjz = (_wuGong + GetArrtValue(ItemInfo.Equip, pb.Attribute.PHYSICAL_ATTACK) + (float)cfgInfo.Phyatk * equipValue) * (1f + GetArrtValue(ItemInfo.Equip, pb.Attribute.PHYSICAL_ATTACK_PERCENTAGE) / baseFM);
            else
                zgjz = (_moGong + GetArrtValue(ItemInfo.Equip, pb.Attribute.MAGIC_ATTACK) + (float)cfgInfo.Magatk * equipValue) * (1f + GetArrtValue(ItemInfo.Equip, pb.Attribute.MAGIC_ATTACK_PERCENTAGE) / baseFM);

            //增伤比
            float zsb = GetArrtValue(ItemInfo.Equip, pb.Attribute.DMT_FINAL_INCREASE_PERCENTAGE) + (float)cfgInfo.PVEdamage * equipValue + (float)cfgInfo.PVPdamage * equipValue + GetArrtValue(ItemInfo.Equip, pb.Attribute.PVP_DAMAGE_PERCENTAGE) + GetArrtValue(ItemInfo.Equip, pb.Attribute.PVE_DAMAGE_PERCENTAGE) + GetArrtValue(ItemInfo.Equip, pb.Attribute.DMT_FINAL_INCREASE_PERCENTAGE);

            //暴击增伤百分比
            float bjzs = _baoJi + GetArrtValue(ItemInfo.Equip, pb.Attribute.CRITICAL_DAMAGE_PERCENTAGE);

            //增伤固定总值
            float zsgdzz = (GetArrtValue(ItemInfo.Equip, pb.Attribute.DMT_INCREASE_FIXED_PLAYER) + GetArrtValue(ItemInfo.Equip, pb.Attribute.DMT_INCREASE_FIXED_MONSTER)) / 2f + GetArrtValue(ItemInfo.Equip, pb.Attribute.DMT_FINAL_INCREASE_FIXED);

            //有效输出
            float yxsc = zgjz * (1 + zsb / baseFM + bjzs / baseFM) * (1 + _skillCombat / 100) + zsgdzz;

            float allPower = (float)Math.Sqrt(AllHP * CFG_XMLConfig.xmlConfig.BaseValue.Power.lifevalue) * (float)Math.Sqrt(yxsc * CFG_XMLConfig.xmlConfig.BaseValue.Power.attackvalue)
                + ((float)cfgInfo.Hit * equipValue + GetArrtValue(ItemInfo.Equip, pb.Attribute.HIT)) * (1 + GetArrtValue(ItemInfo.Equip, pb.Attribute.HIT_PERCENTAGE)) * CFG_XMLConfig.xmlConfig.BaseValue.Power.hitvalue
                + ((float)cfgInfo.Miss * equipValue + GetArrtValue(ItemInfo.Equip, pb.Attribute.DODGE)) * (1 + GetArrtValue(ItemInfo.Equip, pb.Attribute.DODGE_PERCENTAGE)) * CFG_XMLConfig.xmlConfig.BaseValue.Power.dodgevalue
                + ((float)cfgInfo.Crit * equipValue + GetArrtValue(ItemInfo.Equip, pb.Attribute.CRITICAL_STRIKE)) * (1 + GetArrtValue(ItemInfo.Equip, pb.Attribute.CRITICAL_STRIKE_PERCENTAGE)) * CFG_XMLConfig.xmlConfig.BaseValue.Power.critvalue
                + ((float)cfgInfo.Decrit * equipValue + GetArrtValue(ItemInfo.Equip, pb.Attribute.TOUGHNESS)) * (1 + GetArrtValue(ItemInfo.Equip, pb.Attribute.TOUGHNESS_PERCENTAGE)) * CFG_XMLConfig.xmlConfig.BaseValue.Power.toughness;

            //Debug.Log($" 对比槽位 {_TeamIdx} 唯一ID{ItemInfo.Uid}  装备 {cfgInfo.Name} 有效值{AllHP} 生命值 {hp} 物防 {wf} 魔防 {mf} 增伤比 {zsb} 暴击增伤百分比 {bjzs} 增伤固定总值 {zsgdzz} 有效输出 {yxsc} 总战力 {allPower}");

            return (int)Math.Round(allPower);
        }

        /// <summary>
        /// 获取装备战力
        /// </summary>
        /// <param name="soltIdx">玩家槽位 0为玩家 1为伙伴</param>
        /// <param name="ItemInfo"> 装备信息</param>
        /// <returns></returns>
        private (int, int) GetEquiopBsetByOnePartner(int soltIdx, pb.ItemInfo ItemInfo)
        {
            var equipCFG = CFG_Item.GetById(ItemInfo.CfgId);
            if (soltIdx == 0)
            {
                if (equipCFG.Occupylimit != 0 && RoleM.Instance.Job != equipCFG.Occupylimit)
                    return (0, -1);

                if (equipCFG.Levellimit > RoleM.Instance.Level)
                    return (0, -1);

                if (RoleM.Instance.CheckIsNullEquip((pb.EquipSlot)equipCFG.Equiptype))
                    return (1, soltIdx);

                return (GetEquipPower(ItemInfo, RoleM.Instance.Job, soltIdx) - GetEquipPower(RoleM.Instance.GetSlotInfoBySolt((pb.EquipSlot)equipCFG.Equiptype).Item, RoleM.Instance.Job, soltIdx), soltIdx);
            }
            else
            {
                var solt = PartnerM.Instance.GetOnePartnerSlotInfo(soltIdx);
                if (solt.PartnerUid != 0)
                {
                    var partner = PartnerM.Instance.GetOnePartnerInfoByUid(solt.PartnerUid);
                    if (equipCFG.Occupylimit != 0 && (CFG_Partner.GetById(partner.CfgId).job % 10) != equipCFG.Occupylimit)
                    {
                        return (0, -1);
                    }
                    if (partner.Level < equipCFG.Levellimit)
                    {
                        return (0, -1);
                    }
                    var _equipInfo = PartnerM.Instance.GetOnePartnerOneEquip(soltIdx, (int)equipCFG.Equiptype);
                    if (_equipInfo == null)
                        return (1, soltIdx);

                    return (GetEquipPower(ItemInfo, CFG_Partner.GetById(partner.CfgId).job, soltIdx) - GetEquipPower(_equipInfo, CFG_Partner.GetById(partner.CfgId).job, soltIdx), soltIdx);
                }
                else
                {
                    return (0, -1);
                }
            }

        }

        /// <summary>
        /// 检测装备是否更强(第一个值: >0 up、<0 down, 第二个值:=0 主角、>0 伙伴)
        /// </summary>
        /// <param name="ItemInfo"></param>
        /// <returns></returns>
        public (int, int) GetEquipIsBest(pb.ItemInfo ItemInfo, int _TeamIdx = -1)
        {

            if (_TeamIdx != -1)
            {
                return GetEquiopBsetByOnePartner(_TeamIdx, ItemInfo);
            }

            int[] powers = new int[5];
            int[] sords = new int[5];
            //返回一个合适装备战力 +
            for (int i = 0; i <= 4; i++)
            {
                (powers[i], sords[i]) = GetEquiopBsetByOnePartner(i, ItemInfo);
                //返回一个合适装备战力 +
                if (powers[i] > 0)
                    return (powers[i], sords[i]);
            }
            //返回一个合适装备战力 无变动的
            for (int i = 0; i < powers.Length; i++)
            {
                if (powers[i] == 0 && sords[i] >= 0)
                    return (0, sords[i]);
            }
            //返回一个合适装备战力 -
            for (int i = 0; i < powers.Length; i++)
            {
                if (powers[i] < 0)
                    return (powers[i], sords[i]);
            }

            //返回都不合适
            return (0, -1);
        }

        /// <summary>
        /// 装备的极品属性
        /// </summary>
        /// <param name="equipInfo">装备信息</param>
        /// <returns>(属性名，属性值)</returns>
        public (string, string) GetEquipBestAttribute(EquipInfo equipInfo)
        {
            string bestStr1 = "";
            string bestStr2 = "";
            if (equipInfo != null)
            {
                Dictionary<int, int> BestDic = new Dictionary<int, int>();
                foreach (var item in equipInfo.Attributes)
                {
                    BestDic.Add((int)item.attr, item.Num);
                }

                foreach (var item in CFG_XMLConfig.xmlConfig.RolePanelList)
                {
                    int idx = 0;
                    if (BestDic.ContainsKey(item.SID) && BestDic[item.SID] > 0)
                    {
                        int qulite = CFG_RandomDistrib.GetColor((int)equipInfo.Quality, item.SID, BestDic[item.SID]);
                        string arrtStr = CFG_XMLConfig.GetAttributeStrByID(item.SID, BestDic[item.SID]);
                        if (idx % 2 == 0)
                            bestStr1 = $"{bestStr1}{arrtStr}\n";
                        //bestStr1 = $"{bestStr1}{GetStrByQuality(qulite, arrtStr)}\n";
                        else
                            bestStr2 = $"{bestStr2}{arrtStr}\n";
                        //bestStr2 = $"{bestStr2}{GetStrByQuality(qulite, arrtStr)}\n";
                    }
                }
            }
            return (bestStr1, bestStr2);
        }

        public void OnEnegerAdd(Action JumpAction)
        {
            if (RoleM.Instance.Energy >= CFG_XMLConfig.xmlConfig.BaseValue.energy.maxValue)
            {
                UITool.ShowTips("已达到体力值上限！");
                return;

            }

            var itemInfo = CFG_Store.GetById(10001);
            UITool.Open<UIMessageBox>(null, new UIMessageBox.CallBack
            {
                ShowStr = $"是否消耗{itemInfo.Limite.Num}点券购买{itemInfo.Solditems[0].Num}点体力？\n 今日限购次数{itemInfo.Limite.limitnum - StoreM.Instance.GetBuyNumByID(itemInfo.ID)}/{itemInfo.Limite.limitnum}",
                YesFun = () =>
                {
                    if (!ComCoinToCheck(10001, JumpAction))
                    {
                        return;
                    }
                    GameToBuy(10001, null);
                },
                BtnYes = "购买",
            });
        }

        public void SetEffectLayer(Transform effect, int layer)
        {
            var ps = effect.GetComponentsInChildren<Renderer>();
            for (int i = 0; i < ps.Length; i++)
            {
                ps[i].sortingOrder = layer;
            }

            var canvas = effect.GetComponentsInChildren<Canvas>();
            for (int i = 0; i < canvas.Length; i++)
            {
                canvas[i].sortingOrder = layer;
            }
        }

        public bool CheakIsLock(CFG_XMLConfig.XML_Config.SwitchOpne switchOpne, bool showTip = true, StringBuilder strTextTips = default)
        {
            string strtip = "";
            if (CFG_XMLConfig.xmlConfig.SwitchOpenLvDic.ContainsKey(switchOpne) && CFG_XMLConfig.xmlConfig.SwitchOpenLvDic[switchOpne] > RoleM.Instance.Level)
            {
                strtip = CFG_XMLConfig.xmlConfig.SwitchOpenLvDic[switchOpne] + "级解锁";
                strTextTips?.Append(strtip);
                if (showTip)
                    UITool.ShowTips(strtip);
                return true;
            }


            CFG_XMLConfig.XML_Config.XML_OpenFunction curFuncData = null; //当前展示功能数据
            if (CFG_XMLConfig.xmlConfig.SwitchOpenFunction.ContainsKey(switchOpne))
                curFuncData = CFG_XMLConfig.xmlConfig.SwitchOpenFunction[switchOpne];

            if (curFuncData == null) return false;//没有读到配置就就当满足开启条件

            if (curFuncData.openMinBattel > RoleM.Instance.MainBattleId)
            {
                var battleInfo = CFG_MainBattle.GetById(curFuncData.openMinBattel);
                if (battleInfo == null) return false;//没有找到配置信息，就当满足开启条件
                //strtip = $"通关{battleInfo.name}";
                strtip = $"{battleInfo.name.Replace("主线关卡", "主线关卡\n")}解锁";
                //Debug.LogError("=========================="+ RoleM.Instance.MainBattleId);
                strTextTips?.Append(strtip);
                if (showTip)
                {
                    UITool.ShowTips($"通关{battleInfo.name}解锁");
                }

                return true;
            }

            return false;
        }

        private bool mainStageIsPass(CFG_XMLConfig.XML_Config.SwitchOpne curID)
        {
            CFG_XMLConfig.XML_Config.XML_OpenFunction curFuncData = null; //当前展示功能数据
            if (CFG_XMLConfig.xmlConfig.SwitchOpenFunction.ContainsKey(curID))
                curFuncData = CFG_XMLConfig.xmlConfig.SwitchOpenFunction[curID];

            if (curFuncData == null) return false;//没有读到配置就不显示

            return curFuncData.openMinBattel <= RoleM.Instance.MainBattleId;
        }

        /// <summary>
        /// 开服后，功能已到开启时间后是否提示过新功能说明？
        /// </summary>
        /// <param name="openDay"></param>
        /// <returns></returns>
        public bool CheckIsTipNewFunction(int openDay)
        {
            return GetPlayerPrefsInt($"NewFunctionOpenTip_{RoleM.Instance.Uid}_{openDay}") == 1;
        }

        /// <summary>
        /// 开服后，功能已到开启时间后，记录提示过新功能说明。
        /// </summary>
        /// <param name="openDay"></param>
        /// <param name="isTip"></param>
        public void SetIsTipNewFunction(int openDay, bool isTip)
        {
            SetPlayerPrefs($"NewFunctionOpenTip_{RoleM.Instance.Uid}_{openDay}", isTip ? 1 : 0);
        }

        /// <summary>
        /// 开服后，功能已到开启时间后，弹出过新功能说明。
        /// </summary>
        /// <param name="openDay"></param>
        public void ShowNewFunctionTip(int openDay)
        {
            StringBuilder sb = new StringBuilder();
            sb.AppendLine("<size=24>今日开启功能</size>");
            int number = 0;
            foreach (var item in CFG_XMLConfig.xmlConfig.SwitchOpenFunction)
            {
                if (item.Value.openDay == openDay && item.Value.isShow)
                {
                    number++;
                    sb.AppendLine($"{number}.{item.Value.describe}");
                }
            }
            UITool.Open<UINewFunctionOpen>(null, sb.ToString());
        }

        public bool HaveNewFunctionTip(int openDay)
        {
            foreach (var item in CFG_XMLConfig.xmlConfig.SwitchOpenFunction)
            {
                if (item.Value.openDay == openDay && item.Value.isShow)
                    return true;
            }
            return false;
        }

        /// <summary>
        /// 检查开服后功能是否到达开启时间
        /// </summary>
        /// <param name="switchOpne"></param>
        /// <param name="showTip"></param>
        /// <returns></returns>
        public bool CheckOpenServerIsLock(CFG_XMLConfig.XML_Config.SwitchOpne switchOpne, bool showTip = true)
        {
            XML_OpenFunction curFuncData; //当前展示功能数据
            if (CFG_XMLConfig.xmlConfig.SwitchOpenFunction.ContainsKey(switchOpne))
                curFuncData = CFG_XMLConfig.xmlConfig.SwitchOpenFunction[switchOpne];
            else return false;

            bool isLock;
            if(curFuncData.severday == 0)
            {
                isLock = !IsArrivedRegistRoleTime(curFuncData.openDay);
                if (isLock && showTip)
                {
                    int days = DaysRegistRoleTimeToNow(curFuncData.openDay);
                    UITool.ShowTips($"{days}天后开启");
                }
            }
            else
            {
                long openTime = GetServerOpenTimeMillisecond();
                long openZeroTime = ConvertToMidnightSeconds(openTime);
                long opday = (GetCurrentTimeSeconds() - openZeroTime) / 86400;
                int cfgOpenDay = GetOpenServerUnlockDay(switchOpne);
                isLock = opday < cfgOpenDay - 1;
                if (isLock && showTip)
                {
                    UITool.ShowTips($"{cfgOpenDay - opday - 1}天后开启");
                }
            }

            return isLock;
        }

        /// <summary>
        /// 开服的第几天
        /// </summary>
        public int GetTheDayOpenServer()
        {
            long openTime = GetServerOpenTimeMillisecond();
            long openZeroTime = ConvertToMidnightSeconds(openTime);
            long opday = (GetCurrentTimeSeconds() - openZeroTime) / 86400;
            return (int)opday + 1;
        }

        /// <summary>
        /// 开服后，距离功能开启还剩多少天？（开服几天后开启功能）
        /// </summary>
        /// <returns>-1：未到开服时间</returns>
        public int RemainingDaysUnlockOpen(CFG_XMLConfig.XML_Config.SwitchOpne switchOpne)
        {
            XML_OpenFunction curFuncData; //当前展示功能数据
            if (CFG_XMLConfig.xmlConfig.SwitchOpenFunction.ContainsKey(switchOpne))
                curFuncData = CFG_XMLConfig.xmlConfig.SwitchOpenFunction[switchOpne];
            else return 0;

            if (curFuncData.severday == 0)
            {
                return DaysRegistRoleTimeToNow(curFuncData.openDay);
            }
            else
            {
                long openTime = GetServerOpenTimeMillisecond();
                long openZeroTime = ConvertToMidnightSeconds(openTime);
                long opday = (GetCurrentTimeSeconds() - openZeroTime) / 86400;
                int cfgOpenDay = GetOpenServerUnlockDay(switchOpne);
                return cfgOpenDay - (int)opday - 1;
            }
        }

        /// <summary>
        /// 剩余有效天数。
        /// 最后一天显示0天
        /// </summary>
        /// <param name="millisecondsTimestamp">有效期截止时间戳</param>
        /// <param name="validDays">有效天数</param>
        /// <returns></returns>
        public int RemainingDays(long millisecondsTimestamp)
        {
            long endZeroTime = ConvertToMidnightSeconds(millisecondsTimestamp);
            DateTime endTimeData = ConvertFromUnixTimestampSeconds(endZeroTime);

            long todayZeroTime = ConvertToMidnightSeconds(GetCurrentTimeSeconds() * 1000);
            DateTime todayTimeData = ConvertFromUnixTimestampSeconds(todayZeroTime);

            return (endTimeData - todayTimeData).Days;
        }

        /// <summary>
        /// 开服第几天功能开启？
        /// </summary>
        /// <param name="switchOpne"></param>
        /// <returns></returns>
        public int GetOpenServerUnlockDay(CFG_XMLConfig.XML_Config.SwitchOpne switchOpne)
        {
            if (CFG_XMLConfig.xmlConfig.SwitchOpenDayDic.ContainsKey(switchOpne))
                return CFG_XMLConfig.xmlConfig.SwitchOpenDayDic[switchOpne];

            return 1;
        }

        /// <summary>
        /// 卷轴兑换挑战次数
        /// </summary>
        /// <param name="itemId">卷轴道具id</param>
        /// <param name="isJump">卷轴不足时跳转到获取页面</param>
        public void CommonChallengeCountExchange(CFG_Item.ItemID itemId, bool isJump = true)
        {
            CommonChallengeCountExchange((int)itemId, isJump);
        }

        /// <summary>
        /// 卷轴兑换挑战次数
        /// </summary>
        /// <param name="itemId">卷轴道具id</param>
        /// <param name="isJump">卷轴不足时跳转到获取页面</param>
        public void CommonChallengeCountExchange(int itemId, bool isJump = true)
        {
            //string title = itemId switch
            //{
            //    25073 => "首领挑战次数兑换",
            //    25074 => "遗迹挑战次数兑换",
            //    25075 => "强化石副本挑战次数兑换",
            //    25076 => "升星石副本挑战次数兑换",
            //    25077 => "魂石副本挑战次数兑换",
            //    25078 => "书页副本挑战次数兑换",
            //    25079 => "金币副本挑战次数兑换",
            //    25080 => "经验副本挑战次数兑换",
            //    25081 => "五行试炼挑战次数兑换",
            //    25082 => "即时赛挑战次数兑换",
            //    25083 => "回合赛挑战次数兑换",
            //    _ => "",
            //};
            string title = null;
            string desc2 = null;
            if (itemId == 56)
            {
                title = "挂机时间增加提示";
                var cfgitem = CFG_Item.GetById(itemId);
                long canusenum = CFG_ItemUse.GetDailyUseLimit(itemId) - BagM.Instance.GetDailyUsedItem(itemId);
                desc2 = $"（今日还可使用：{canusenum}）";
            }
            else if (itemId >= 25073 && itemId <= 25083)
            {
                title = "次数增加提示";
            }

            UITool.Open<UICommonItemBuy>(null, new UICommonItemBuy.BuyInfo
            {
                option = UICommonItemBuy.Option.Use,
                title = title,
                itemId = itemId,
                jumpShopTickets = isJump,
                description2 = desc2,
            });
        }


        // 槽位是否解锁是否有角色
        public bool SlotIsUnLockOrHasRole(int slotIdx)
        {
            AudioManager.Instance.PlaySound(AudioKey.Sound.Button);
            if (slotIdx <= 1)
                return true; // 主角、随机槽位永远开启

            if (slotIdx > 1 && CFG_XMLConfig.CheckPartnerSoltIsOpenServer(slotIdx - 1))
            {
                UITool.ShowTips($"槽位尚未解锁，无法选择！");
                return false;
            }

            if (CFG_XMLConfig.GetOpenSoltLv(slotIdx - 1) > RoleM.Instance.Level)
            {
                UITool.ShowTips("槽位尚未解锁，无法选择！");
                return false;
            }

            var solt = PartnerM.Instance.GetOnePartnerSlotInfo(slotIdx - 1);
            if (solt.PartnerUid == 0)
                UITool.ShowTips("槽位为空，无法选择！");

            return solt.PartnerUid != 0;
        }

        /// <summary>
        /// 装备槽位属性解锁状态
        /// </summary>
        /// <param name="slot"></param>
        /// <returns></returns>
        public bool EquipAttributeSlotIsUnlock(CharacterSlot characterSlot, EquipSlot equipSlot, RefinedSlot refinedSlot)
        {
            int attrSlot = (int)refinedSlot;

            (int type, int limit) = CFG_XMLConfig.GetOpenEquipSlotAttribute(attrSlot);

            if (type == 1)
                return RoleM.Instance.Level >= limit;
            else
                return RoleM.Instance.IsRefineSlotBuy(characterSlot, equipSlot, refinedSlot);
        }

        public string GetEquipAttributeSlotIsUnlockDesc(RefinedSlot refinedSlot)
        {
            var lockcfg = CFG_XMLConfig.xmlConfig.GearRefiningDic[(int)refinedSlot];
            if (lockcfg.type == 1)
                return $"角色{lockcfg.limit}级解锁";
            else if (lockcfg.type == 2)
                return $"￥{lockcfg.limit}解锁";
            else if (lockcfg.type == 3)
                return $"       {lockcfg.limit}";
            else
                return "";
        }

        public Tweener MoveScrollViewContentByOffset(LoopListView2 loopListView, float offset, float moveDuration = 0, Ease ease = Ease.InOutQuad, TweenCallback onMoveCompleted = null)
        {
            if (Object.Equals(null, loopListView) || offset == 0 || loopListView.ItemTotalCount == 0)
            {
                return null;
            }

            ScrollRect scrollRect = loopListView.GetComponent<ScrollRect>();

            if (loopListView.ArrangeType == ListItemArrangeType.TopToBottom)
            {
                Vector3 pos = scrollRect.content.anchoredPosition3D;
                pos.y += offset;
                return scrollRect.content.DOAnchorPos3D(pos, moveDuration).SetEase(ease).OnComplete(onMoveCompleted);
            }
            else if (loopListView.ArrangeType == ListItemArrangeType.BottomToTop)
            {
                Vector3 pos = scrollRect.content.anchoredPosition3D;
                pos.y += offset;
                return scrollRect.content.DOAnchorPos3D(pos, moveDuration).SetEase(ease).OnComplete(onMoveCompleted);
            }
            else if (loopListView.ArrangeType == ListItemArrangeType.LeftToRight)
            {
                Vector3 pos = scrollRect.content.anchoredPosition3D;
                pos.x += offset;
                return scrollRect.content.DOAnchorPos3D(pos, moveDuration).SetEase(ease).OnComplete(onMoveCompleted);
            }
            else if (loopListView.ArrangeType == ListItemArrangeType.RightToLeft)
            {
                Vector3 pos = scrollRect.content.anchoredPosition3D;
                pos.x += offset;
                return scrollRect.content.DOAnchorPos3D(pos, moveDuration).SetEase(ease).OnComplete(onMoveCompleted);
            }
            return null;
        }

        #region 游戏设置项保存
        public void SetPlayerPrefs(string key, string value)
        {
            PlayerPrefsExtend.SetString(key, value);
        }

        public string GetPlayerPrefs(string key)
        {
            return GetPlayerPrefs(key, "");
        }

        public string GetPlayerPrefs(string key, string defaultValue)
        {
            return PlayerPrefsExtend.GetString(key, defaultValue);
        }

        public void SetPlayerPrefs(string key, int value)
        {
            SetPlayerPrefs(key, value.ToString());
        }

        public int GetPlayerPrefsInt(string key)
        {
            if (int.TryParse(GetPlayerPrefs(key), out int value))
            {
                return value;
            }
            return 0;
        }

        public int GetPlayerPrefsInt(string key, int defaultValue)
        {
            if (int.TryParse(GetPlayerPrefs(key), out int value))
            {
                return value;
            }
            return defaultValue;
        }

        public void SetPlayerPrefs(string key, float value)
        {
            SetPlayerPrefs(key, value.ToString());
        }

        public float GetPlayerPrefsFloat(string key)
        {
            if (float.TryParse(GetPlayerPrefs(key), out float value))
            {
                return value;
            }
            return 0;
        }

        public float GetPlayerPrefsFloat(string key, float defaultValue)
        {
            if (float.TryParse(GetPlayerPrefs(key), out float value))
            {
                return value;
            }
            return defaultValue;
        }

        public void SetPlayerPrefs(string key, long value)
        {
            SetPlayerPrefs(key, value.ToString());
        }

        public long GetPlayerPrefsLong(string key)
        {
            if (long.TryParse(GetPlayerPrefs(key), out long value))
            {
                return value;
            }
            return 0;
        }
        public long GetPlayerPrefsLong(string key, long defaultValue)
        {
            if (long.TryParse(GetPlayerPrefs(key), out long value))
            {
                return value;
            }
            return defaultValue;
        }

        public void DeletePlayerPrefs(string key)
        {
            PlayerPrefsExtend.DeleteKey(key);
        }

        #endregion

        #region 开服后功能是否开启
        // 获取当前服务器名称
        public string GetServerName()
        {
            var selectSever = GameM.Instance.SeverInfoDic[GameM.Instance.SeverDicIdx][GameM.Instance.SeverIdx];
            return selectSever.name;
        }

        /// <summary>
        /// 开服时间
        /// </summary>
        /// <returns></returns>
        public long GetServerOpenTime()
        {
            var selectSever = GameM.Instance.SeverInfoDic[GameM.Instance.SeverDicIdx][GameM.Instance.SeverIdx];
            return selectSever.open_by / 1000;
        }

        /// <summary>
        /// 开服时间
        /// </summary>
        /// <returns></returns>
        public long GetServerOpenTimeMillisecond()
        {
            var selectSever = GameM.Instance.SeverInfoDic[GameM.Instance.SeverDicIdx][GameM.Instance.SeverIdx];
            return selectSever.open_by;
        }

        /// <summary>
        /// 已开服时间
        /// </summary>
        /// <returns></returns>
        public long GetServerOpeningTime()
        {
            return GetCurrentTimeSeconds() - GetServerOpenTime();
        }

        #endregion

        public string GetJobName(int jobId)
        {
            if (jobId % 10 == 1)
                return "战士";
            else if (jobId % 10 == 2)
                return "法师";
            else if (jobId % 10 == 3)
                return "道士";
            else return "";
        }

        /// <summary>
        /// 设置头像
        /// </summary>
        /// <param name="cfgId">配置id</param>
        /// <param name="job">职业</param>
        /// <param name="imgHead">头像</param>
        /// <param name="imgQuality">头像品质</param>
        public void SetRoleHead(int cfgId, int job, Image imgHead, Image imgQuality = null)
        {
            var cfgHead = CFG_RoleIcon.GetById(cfgId);
            if (!Equals(imgHead, null))
                ResouseHelper.Instance.LoadSpriteByAtlas(cfgHead.Path, cfgHead.GetIcon(job), imgHead);
            if (!Equals(imgQuality, null))
                ResouseHelper.Instance.LoadSpriteByAtlas(cfgHead.Path, GameConst.HeadIconQualityPrefix + cfgHead.Quality, imgQuality);
        }

        /// <summary>
        /// 设置头像框
        /// </summary>
        /// <param name="cfgId">配置id</param>
        /// <param name="job">职业</param>
        /// <param name="imgFrame">头像框</param>
        public void SetRoleHeadFrame(int cfgId, int job, Image imgFrame)
        {
            var cfgFrame = CFG_RoleIcon.GetById(cfgId);
            if (!Equals(imgFrame, null))
                ResouseHelper.Instance.LoadSpriteByAtlas(cfgFrame.Path, cfgFrame.GetIcon(job), imgFrame);
        }

        public void SetDefauleRoleHead(int job, Image imgFrame, Image imgQuality)
        {
            var cfgFrame = CFG_RoleIcon.GetById(CFG_XMLConfig.xmlConfig.BaseValue.AvatarFrame);
            if (!Equals(imgFrame, null))
                ResouseHelper.Instance.LoadSpriteByAtlas(cfgFrame.Path, cfgFrame.GetIcon(job), imgFrame);
            if (!Equals(imgQuality, null))
                ResouseHelper.Instance.LoadSpriteByAtlas(cfgFrame.Path, GameConst.HeadIconQualityPrefix + 1, imgQuality);
        }

        /// <summary>
        /// 获取正确的头像和头像框id。
        /// 如果传入0则返回默认头像和头像框id。
        /// </summary>
        /// <param name="headId">头像id</param>
        /// <param name="frameId">头像框id</param>
        /// <returns></returns>
        public (int, int) GetHeadAndFrameId(int headId, int frameId)
        {
            headId = headId == 0 ? CFG_XMLConfig.xmlConfig.BaseValue.Avatar : headId;
            frameId = frameId == 0 ? CFG_XMLConfig.xmlConfig.BaseValue.AvatarFrame : frameId;
            return (headId, frameId);
        }

        #region 数数上报
        public void SendTDGuideCompleted(int guideStep, string guideName)
        {
            Dictionary<string, object> properties = new Dictionary<string, object>();
            properties["guide_step"] = guideStep;
            properties["guide_name"] = guideName;
            TDAnalyticsUtil.Track("guide_completed", properties, true);
        }

        // 数数公共事件设置
        public void TDPublicEvent()
        {
            Dictionary<string, object> superProperties = new Dictionary<string, object>();
            superProperties["channel"] = "微信";
            superProperties["server_id"] = GameM.Instance.SeverInfoDic[GameM.Instance.SeverDicIdx][GameM.Instance.SeverIdx].server_id;
            superProperties["nick_name"] = RoleM.Instance.Nickname;
            superProperties["tmp_power"] = RoleM.Instance.TotalCombat;
            superProperties["tmp_level"] = RoleM.Instance.Level;
            superProperties["tmp_box_level"] = RoleM.Instance.BoxNowlevel;
            superProperties["tmp_main_id"] = GetZhanYiLevel();
            superProperties["tmp_task_id"] = GetTaskAndProgress();
            superProperties["tmp_role_job"] = GetJobName(RoleM.Instance.Job);
            superProperties["tmp_account"] = GameM.Instance.WXOpneID;
            superProperties["tmp_device"] = "";
            superProperties["tmp_role_id"] = RoleM.Instance.Uid;

            var petinfo = PetM.Instance.GetOnePetByIdx(0);
            bool isZun = false;

            superProperties["tmp_pet_id"] = petinfo.CfgId;
            superProperties["tmp_pet_name"] = petinfo.CfgId == 0 ? "" : CFG_Pet.GetById(petinfo.CfgId).Name;
            superProperties["tmp_pet_level"] = petinfo.Level;
            superProperties["tmp_pet_rank"] = petinfo.Star;

            if (petinfo.CfgId != 0)
            {
                var cfgInfo = CFG_Pet.GetById(petinfo.CfgId);
                isZun = cfgInfo.Mark == 1;
            }
            superProperties["tmp_pet_zun"] = isZun;

            TDAnalyticsUtil.SetSuperProperties(superProperties);
            Debug.Log("公共事件设置成功！");
        }

        /// <summary>
        /// 战役关卡进度
        /// </summary>
        /// <returns></returns>
        public string GetZhanYiLevel()
        {
            //单机服：副本信息可能还没请求回来（GetInfo() 为 null），埋点里不能崩
            var info = FightPlayM.Instance.GetInfo();
            if (info == null || info.battle_map == null)
                return "关卡1-1";
            var fightInfo = info.battle_map;
            var CFGBattle = CFG_Battle.GetNextFightById(fightInfo.BattleId);
            if (CFGBattle == null)
                return "关卡1-1";
            return "关卡" + CFGBattle.Chapter + "-" + CFGBattle.Layers;
        }

        public string GetTaskAndProgress()
        {
            var tasks = TaskM.Instance.GetTasksByClass(TaskClass.MAINLINE);
            if (tasks == null || tasks.Count == 0)
                return "暂无主线任务";
            var task = tasks[0];
            (int amout, int targetnum) = CFG_Task.GetTaskPro(task.CfgId, task.Target.Amount, task.Target.TargetNum);
            return $"{CFG_Task.StrTaskDes(task.CfgId)}({FormatNumWithUnit(amout)}/{FormatNumWithUnit(targetnum)})";
        }

        #endregion

        #region wx tool
        public string GetWXLaunchParamClickId()
        {
            var launchParams = SDKCtr.Instance.GetWXLaunchParams();
            /*
            Debug.Log($"微信启动参数：GetWXLaunchParams(), weixin:外部url = ");
            JsonMapper.ToJson( _param );
            string paramkey = null;
            if (_param != null)
            {
                foreach (var item in _param)
                {
                    Debug.Log("_param key = " + item.Key + " |\n (UnEscape) " + UnityWebRequest.UnEscapeURL(item.Key));
                    Debug.Log("_param Value = " + item.Value + " |\n (UnEscape) " + UnityWebRequest.UnEscapeURL(item.Value));
                    if (item.Key == "clickid")
                    {
                        Debug.Log("------ wx launch clickid = " + item.Value);
                        return item.Value;
                    }
                    else if (item.Key.Contains("clickid"))
                    {
                        paramkey = item.Key;
                        Debug.Log($"找到 clickid: key = {item.Key}; value = {item.Value}");
                        break;
                    }
                }
            }

            string clickid = "";

            if (!string.IsNullOrEmpty(paramkey))
            {
                var kvarray = UnityWebRequest.UnEscapeURL(paramkey).Split('&');
                foreach (var kv in kvarray)
                {
                    var kvpair = kv.Split('=');
                    if (kvpair.Length == 2 && kvpair[0] == "clickid")
                    {
                        clickid = kvpair[1];
                        break;
                    }
                }
            }
            */
            string value = JsonMapper.ToJson(launchParams);
            Debug.Log("------ wx launch param = " + value);
            return value;
        }
        #endregion
    }
}
