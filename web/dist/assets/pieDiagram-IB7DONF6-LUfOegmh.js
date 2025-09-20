import{p as U}from"./chunk-4BMEZGHF-CEzR6VGu.js";import{ac as y,a4 as z,aH as j,E as H,o as K,p as Q,b as Z,g as q,d as J,c as X,_ as p,l as F,v as Y,e as tt,F as et,K as at,Q as rt,m as nt}from"./index-D8gcSXCW.js";import{p as it}from"./radar-MK3ICKWK-2Aqjm5Y1.js";import{d as P}from"./arc-Dn5axUjj.js";import{o as st}from"./ordinal-Cboi1Yqb.js";import"./semi-ui-ClUhRvP3.js";import"./react-core-DtMAePju.js";import"./tools-CmMbhh-g.js";import"./react-components-XvIrFW7f.js";import"./i18n-CtMfBwdO.js";import"./_baseUniq-DXBUCMUV.js";import"./_basePickBy-BPJYDTjx.js";import"./clone-D3UoTbVK.js";import"./init-Gi6I4Gst.js";function ot(t,a){return a<t?-1:a>t?1:a>=t?0:NaN}function lt(t){return t}function ct(){var t=lt,a=ot,m=null,o=y(0),u=y(z),x=y(0);function i(e){var r,l=(e=j(e)).length,g,A,h=0,c=new Array(l),n=new Array(l),v=+o.apply(this,arguments),w=Math.min(z,Math.max(-z,u.apply(this,arguments)-v)),f,T=Math.min(Math.abs(w)/l,x.apply(this,arguments)),$=T*(w<0?-1:1),d;for(r=0;r<l;++r)(d=n[c[r]=r]=+t(e[r],r,e))>0&&(h+=d);for(a!=null?c.sort(function(S,C){return a(n[S],n[C])}):m!=null&&c.sort(function(S,C){return m(e[S],e[C])}),r=0,A=h?(w-l*$)/h:0;r<l;++r,v=f)g=c[r],d=n[g],f=v+(d>0?d*A:0)+$,n[g]={data:e[g],index:r,value:d,startAngle:v,endAngle:f,padAngle:T};return n}return i.value=function(e){return arguments.length?(t=typeof e=="function"?e:y(+e),i):t},i.sortValues=function(e){return arguments.length?(a=e,m=null,i):a},i.sort=function(e){return arguments.length?(m=e,a=null,i):m},i.startAngle=function(e){return arguments.length?(o=typeof e=="function"?e:y(+e),i):o},i.endAngle=function(e){return arguments.length?(u=typeof e=="function"?e:y(+e),i):u},i.padAngle=function(e){return arguments.length?(x=typeof e=="function"?e:y(+e),i):x},i}var R=H.pie,G={sections:new Map,showData:!1,config:R},b=G.sections,W=G.showData,pt=structuredClone(R),ut=p(()=>structuredClone(pt),"getConfig"),gt=p(()=>{b=new Map,W=G.showData,Y()},"clear"),dt=p(({label:t,value:a})=>{b.has(t)||(b.set(t,a),F.debug(`added new section: ${t}, with value: ${a}`))},"addSection"),ft=p(()=>b,"getSections"),mt=p(t=>{W=t},"setShowData"),ht=p(()=>W,"getShowData"),I={getConfig:ut,clear:gt,setDiagramTitle:K,getDiagramTitle:Q,setAccTitle:Z,getAccTitle:q,setAccDescription:J,getAccDescription:X,addSection:dt,getSections:ft,setShowData:mt,getShowData:ht},vt=p((t,a)=>{U(t,a),a.setShowData(t.showData),t.sections.map(a.addSection)},"populateDb"),St={parse:p(async t=>{const a=await it("pie",t);F.debug(a),vt(a,I)},"parse")},yt=p(t=>`
  .pieCircle{
    stroke: ${t.pieStrokeColor};
    stroke-width : ${t.pieStrokeWidth};
    opacity : ${t.pieOpacity};
  }
  .pieOuterCircle{
    stroke: ${t.pieOuterStrokeColor};
    stroke-width: ${t.pieOuterStrokeWidth};
    fill: none;
  }
  .pieTitleText {
    text-anchor: middle;
    font-size: ${t.pieTitleTextSize};
    fill: ${t.pieTitleTextColor};
    font-family: ${t.fontFamily};
  }
  .slice {
    font-family: ${t.fontFamily};
    fill: ${t.pieSectionTextColor};
    font-size:${t.pieSectionTextSize};
    // fill: white;
  }
  .legend text {
    fill: ${t.pieLegendTextColor};
    font-family: ${t.fontFamily};
    font-size: ${t.pieLegendTextSize};
  }
`,"getStyles"),xt=yt,At=p(t=>{const a=[...t.entries()].map(o=>({label:o[0],value:o[1]})).sort((o,u)=>u.value-o.value);return ct().value(o=>o.value)(a)},"createPieArcs"),wt=p((t,a,m,o)=>{F.debug(`rendering pie chart
`+t);const u=o.db,x=tt(),i=et(u.getConfig(),x.pie),e=40,r=18,l=4,g=450,A=g,h=at(a),c=h.append("g");c.attr("transform","translate("+A/2+","+g/2+")");const{themeVariables:n}=x;let[v]=rt(n.pieOuterStrokeWidth);v??(v=2);const w=i.textPosition,f=Math.min(A,g)/2-e,T=P().innerRadius(0).outerRadius(f),$=P().innerRadius(f*w).outerRadius(f*w);c.append("circle").attr("cx",0).attr("cy",0).attr("r",f+v/2).attr("class","pieOuterCircle");const d=u.getSections(),S=At(d),C=[n.pie1,n.pie2,n.pie3,n.pie4,n.pie5,n.pie6,n.pie7,n.pie8,n.pie9,n.pie10,n.pie11,n.pie12],D=st(C);c.selectAll("mySlices").data(S).enter().append("path").attr("d",T).attr("fill",s=>D(s.data.label)).attr("class","pieCircle");let N=0;d.forEach(s=>{N+=s}),c.selectAll("mySlices").data(S).enter().append("text").text(s=>(s.data.value/N*100).toFixed(0)+"%").attr("transform",s=>"translate("+$.centroid(s)+")").style("text-anchor","middle").attr("class","slice"),c.append("text").text(u.getDiagramTitle()).attr("x",0).attr("y",-400/2).attr("class","pieTitleText");const M=c.selectAll(".legend").data(D.domain()).enter().append("g").attr("class","legend").attr("transform",(s,E)=>{const k=r+l,_=k*D.domain().length/2,B=12*r,V=E*k-_;return"translate("+B+","+V+")"});M.append("rect").attr("width",r).attr("height",r).style("fill",D).style("stroke",D),M.data(S).append("text").attr("x",r+l).attr("y",r-l).text(s=>{const{label:E,value:k}=s.data;return u.getShowData()?`${E} [${k}]`:E});const L=Math.max(...M.selectAll("text").nodes().map(s=>(s==null?void 0:s.getBoundingClientRect().width)??0)),O=A+e+r+l+L;h.attr("viewBox",`0 0 ${O} ${g}`),nt(h,g,O,i.useMaxWidth)},"draw"),Ct={draw:wt},Rt={parser:St,db:I,renderer:Ct,styles:xt};export{Rt as diagram};
