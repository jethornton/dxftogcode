#!/usr/bin/env python

import gtk
import os
import subprocess

class Buglump:

	def __init__(self):
		self.builder = gtk.Builder()
		self.builder.add_from_file('dxf.glade')
		self.builder.connect_signals(self)
		self.window = self.builder.get_object('main_window')
		self.aboutdialog = self.builder.get_object("aboutdialog")
		self.statusbar = self.builder.get_object("statusbar")
		self.context_id = self.statusbar.get_context_id("status")
		self.statusbar.push(self.context_id, 'No File Open')
		self.current_folder = os.path.expanduser('~')
		self.file_name = ''
		self.window.show()

	def on_window_destroy(self, object, data=None):
		gtk.main_quit()

	def on_file_quit(self, menuitem, data=None):
		gtk.main_quit()

	def on_file_open(self, menuitem, data=None):
		self.fcd = gtk.FileChooserDialog("Open...", None,
			gtk.FILE_CHOOSER_ACTION_OPEN,
			(gtk.STOCK_CANCEL, gtk.RESPONSE_CANCEL, gtk.STOCK_OPEN, gtk.RESPONSE_OK))
		print self.fcd.list_shortcut_folders()
		if len(self.current_folder) > 0:
			self.fcd.set_current_folder(self.current_folder)
		self.response = self.fcd.run()
		if self.response == gtk.RESPONSE_OK:
			self.file_name = self.fcd.get_filename()
			self.current_folder = os.path.dirname(self.fcd.get_uri()[7:])
			self.statusbar.push(self.context_id, 'File Selected %s' % str(self.file_name))
		self.fcd.destroy()

	def on_file_convert(self, file_name, data=None):
		if len(self.file_name) > 0:
			print self.file_name
			self.args = self.file_name
			self.result = subprocess.call('dxf2gcode %s' %self.args, shell=True)
			if self.result == 0:
				self.statusbar.push(self.context_id, 'Processing Complete')
			else:
				self.statusbar.push(self.context_id, 'Error Processing %s' % str(self.file_name))
		else:
			self.statusbar.push(self.context_id, 'No File Open')

	def on_view_test(self, item, data=None):
		print subprocess.call('dxf2gcode', shell=True)

	def on_help_about(self, menuitem, data=None):
		self.response = self.aboutdialog.run()
		self.aboutdialog.hide()

if __name__ == '__main__':
  main = Buglump()
  gtk.main()
