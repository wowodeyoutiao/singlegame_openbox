using JEngine.Singleton;
using System.Collections.Generic;
using System.Linq;

namespace HotUpdateScripts.ProtoM
{
    /// <summary>
    /// 邮件信息
    /// </summary>
    public class MailM : SingletonTemplate<MailM>
    {
        private Dictionary<string, HTTP_MailInfo> MailDic = new Dictionary<string, HTTP_MailInfo>();

        public void SetMailLsit(List<HTTP_MailInfo> mailLsit)
        {
            IsNew = false;
            MailDic.Clear();
            foreach (var item in mailLsit)
            {
                MailDic.Add(item.mail_id, item);
            }

        }

        private bool IsNew = false;
        public void SetMailNew(bool _isNew)
        {
            IsNew = _isNew;
        }

        public int GetDotNum()
        {
            if (IsNew)
            {
                return 1;
            }

            List<HTTP_MailInfo> mailList = MailDic.Values.ToList();
            for (int i = 0; i < mailList.Count; i++)
            {
                if (mailList[i].is_read == 0)
                {
                    return 1;
                }
                if (mailList[i].is_attachments == 1 && mailList[i].is_extract == 0)
                {
                    return 1;
                }
            }
            return 0;
        }

        public List<HTTP_MailInfo> GetMailList()
        {
            List<HTTP_MailInfo> mailList = MailDic.Values.ToList();
            mailList.Sort((a, b) =>
            {
                if (a.created_at > b.created_at)
                {
                    return -1;
                }
                else
                {
                    return 1;
                }
            });
            return mailList;
        }

        public void SetOneMailInfo(HTTP_MailInfo mailInfo)
        {
            MailDic[mailInfo.mail_id] = mailInfo;
        }

        /// <summary>
        /// 设置领取状态
        /// </summary>
        /// <param name="mainId"></param>
        public void SetMainDelByID(string mainId)
        {
            if (MailDic.ContainsKey(mainId))
            {
                MailDic.Remove(mainId);
            }

        }

    }


    /// <summary>
    /// 请求邮件信息买返回
    /// </summary>
    public struct HTTP_MailListResponse
    {
        public int code { get; set; }

        //总数量
        public int total { get; set; }

        //列表
        public List<HTTP_MailInfo> data { get; set; }

    }

    public struct HTTP_MailInfo
    {
        public string mail_id { get; set; }

        public string title { get; set; }

        public string desc { get; set; }

        public int is_read { get; set; }

        /// <summary>
        /// 是否提取 0 no
        /// </summary>
        public int is_extract { get; set; }

        /// <summary>
        /// 是否有附件 0 no
        /// </summary>
        public int is_attachments { get; set; }

        /// <summary>
        ///  创建时间
        /// </summary>
        public long created_at { get; set; }

        /// <summary>
        /// 截至时间
        /// </summary>
        public long expiry_at { get; set; }

        //附件
        public Http_AttachmentInfo[] attachments { get; set; }
    }

    /// <summary>
    /// 邮件附件信息
    /// </summary>
    public struct Http_AttachmentInfo
    {
        public int user_props_id { get; set; }
        public int num { get; set; }
    }

    /// <summary>
    /// 邮件详情
    /// </summary>
    public struct Http_MailDetailRequset
    {
        public string mail_id { get; set; }
    }

    /// <summary>
    /// 邮件详情 返回
    /// </summary>
    public struct Http_MailDetailResponse
    {
        public int code { get; set; }
        public HTTP_MailInfo data { get; set; }
    }

    /// <summary>
    /// 邮件提取 删除 
    /// </summary>
    public struct Http_MailOperateRequset
    {
        public string mail_id { get; set; }
    }

    /// <summary>
    /// 邮件提取 删除 （返回）
    /// </summary>
    public struct Http_MailOperateResponse
    {
        public int code { get; set; }
        public string data { get; set; }
    }

}

