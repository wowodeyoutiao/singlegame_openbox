using System;
using UnityEngine;

namespace MetaFramework.UI
{
    public class RedPointSystems
    {
        /// <summary>
        /// 红点树Root节点
        /// </summary>        
        private RedPointNode rootNode;

        /// <summary>
        /// 示例
        /// </summary>
        public void Example()
        {
            RedPointNode mainNode = new RedPointNode("Main");

            var mailNode = mainNode.AddChildNode("Mail");
            var taskNode = mainNode.AddChildNode("Task");
            var bagNode = mainNode.AddChildNode("Bag");

            mailNode.AddChildNode("Mail1");
            mailNode.AddChildNode("Mail2");
            mailNode.AddChildNode("Mail3");
        }

        /// <summary>
        /// 初始化红点树
        /// </summary>
        public void Init(RedPointNode rootNode)
        {
            this.rootNode = rootNode;
        }

        /// <summary>
        /// 注册红点改变事件
        /// </summary>
        /// <param name="nodeName"></param>
        /// <param name="callBack"></param>
        public void Register(string nodeName, Action<RedPointNode> callBack)
        {
            var node = FindNode(nodeName);
            if (node == null)
            {
                Debug.LogError("register failed! can not find the node:" + nodeName);
                return;
            }
            node.onChangeCount = callBack;
            if (callBack != null)
            {
                callBack(node);
            }
        }

        /// <summary>
        /// 派发事件
        /// </summary>
        /// <param name="nodeName"></param>
        /// <param name="newCount"></param>        
        public void Dispatch(string nodeName, int newCount)
        {
            var node = FindNode(nodeName);
            if (node == null)
            {
                Debug.LogError("register failed! can not find the node:" + nodeName);
                return;
            }
            node.SetRedPointCount(newCount);
            //GameHelper.Instance.SetDeayFrame(() => { node.SetRedPointCount(newCount); });

        }

        public void Dispatch(RedPointConst nodeName)
        {
            if (!RedPointC.RedPointNumDic.ContainsKey(nodeName))
            {
                Debug.LogError("register failed! can not find the RedPointNumDic:" + nodeName);
                return;
            }
            //GameHelper.Instance.SetDeayFrame(() => { Dispatch(nodeName.ToString(), RedPointC.RedPointNumDic[nodeName].Invoke()); });
            Dispatch(nodeName.ToString(), RedPointC.RedPointNumDic[nodeName].Invoke());
        }

        private RedPointNode FindNode(string findNodeName)
        {
            return FindNode(rootNode, findNodeName);
        }

        private RedPointNode FindNode(RedPointNode node, string findNodeName)
        {
            if (node == null || node.dicChilds == null) return null;

            foreach (var item in node.dicChilds)
            {
                if (item.Key == findNodeName)
                    return item.Value;

                var findResult = FindNode(item.Value, findNodeName);
                if (findResult != null)
                    return findResult;
            }
            return null;

            //var arrChilds = node.dicChilds.Keys.ToArray();

            //for (int i = 0; i < arrChilds.Length; i++)
            //{
            //    if (arrChilds[i] == findNodeName)
            //        return node.dicChilds[findNodeName];

            //    var findResult = FindNode(node.dicChilds[arrChilds[i]], findNodeName);

            //    if (findResult != null)
            //        return findResult;
            //}
            //return null;
        }

        public void ClearNodeAndChildCount(RedPointConst nodeName)
        {
            ClearNodeAndChildCount(nodeName.ToString());
        }

        public void ClearNodeAndChildCount(string nodeName)
        {
            var node = FindNode(nodeName);
            ClearNodeCount(node);
        }

        private void ClearNodeCount(RedPointNode node)
        {
            if (node != null)
            {
                node.pointCount = 0;
                if (node.dicChilds.Count > 0)
                {
                    foreach (var child in node.dicChilds.Values)
                    {
                        ClearNodeCount(child);
                    }
                }
            }
        }

    }
}
