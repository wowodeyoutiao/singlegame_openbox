using System;
using System.Collections;
using System.Collections.Generic;

using ThinkingData.Analytics;

using UnityEngine;

public class TDAnalyticsUtil
{
    public static void Login(string account)
    {
        TDAnalytics.Login(account);
    }

    public static void SetDistinctId(string distinctId)
    {
        TDAnalytics.SetDistinctId(distinctId);
    }

    public static string GetDistinctId()
    {
        return TDAnalytics.GetDistinctId();
    }

    public static void SetSuperProperties(Dictionary<string, object> properties)
    {
        TDAnalytics.SetSuperProperties(properties);
    }

    public static void Logout()
    {
        TDAnalytics.Logout();
    }

    public static void Track(string eventName, Dictionary<string, object> properties, bool needFlush = false)
    {
        TDAnalytics.Track(eventName, properties);
        if (needFlush)
        {
            TDAnalytics.Flush();
        }
    }

    public static void Flush()
    {
        TDAnalytics.Flush();
    }

    public static void UserSet(Dictionary<string, object> properties)
    {
        TDAnalytics.UserSet(properties);
    }

    public static void UserSetOnce(Dictionary<string, object> properties)
    {
        TDAnalytics.UserSetOnce(properties);
    }

    public static void UserAdd(Dictionary<string, object> properties)
    {
        TDAnalytics.UserAdd(properties);
    }

    public static void UserAppend(Dictionary<string, object> properties)
    {
        TDAnalytics.UserAppend(properties);
    }

    public static void UserDelete()
    {
        TDAnalytics.UserDelete();
    }

}
