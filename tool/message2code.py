# -*- coding: utf-8 -*-
#
# $ curl https://raw.githubusercontent.com/spotbugs/spotbugs/release-3.1/spotbugs/etc/messages.xml | python message2code.py -n BugDescriptionEn
# $ curl https://raw.githubusercontent.com/spotbugs/spotbugs/release-3.1/spotbugs/etc/messages_ja.xml | python message2code.py -n BugDescriptionJa
# $ curl https://raw.githubusercontent.com/spotbugs/spotbugs/release-3.1/spotbugs/etc/messages_fr.xml | python message2code.py -n BugDescriptionFr

import xml.etree.ElementTree as ElementTree
import sys, re, codecs
from argparse import ArgumentParser

sys.stdout = codecs.lookup('utf_8')[-1](sys.stdout)

parser = ArgumentParser()

parser.add_argument('-n', '--name', dest = 'name', type = str, default = 'BugDescriptionEn')

args = parser.parse_args()

tree = ElementTree.parse(sys.stdin)

print "package main\n"
print "var %s = map[string]string{" % (args.name)

for e in tree.findall('.//BugPattern'):
  code = e.get('type')
  desc = ("%s\n\n%s" % (e.findtext('./ShortDescription'), re.sub(r'\n+', "\n", e.findtext('./Details')).strip())).replace('"', '&apos;').replace("\n", '\\n')

  print '  "%s":"%s",' % (code, desc)

print "}"
