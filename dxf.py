#!/usr/bin/env python

version = '1.0.0'

# Copyright John Thornton 2015

import gtk
import os
#import subprocess
from subprocess import Popen, PIPE
import ConfigParser

class Buglump:

	def __init__(self):
		self.builder = gtk.Builder()
		self.builder.add_from_file('dxf.glade')
		self.builder.connect_signals(self)
		self.window = self.builder.get_object('main_window')
		self.aboutdialog = self.builder.get_object('aboutdialog')
		self.aboutdialog.set_version(version)
		self.file_analyze = self.builder.get_object('file_analyze')
		self.file_convert = self.builder.get_object('file_convert')
		self.tolerance = self.builder.get_object('tolerance_entry')
		self.path = self.builder.get_object('path_entry')
		self.status = self.builder.get_object("status_label")
		self.status.set_text('No File Open')
		self.current_folder = os.path.expanduser('~')
		self.label2 = self.builder.get_object('label2')
		self.save_directory = self.builder.get_object('save_directory_entry')
		self.output_name = self.builder.get_object('output_name_entry')
		self.input_file = ''
		self.ini_file = ''
		self.config = ConfigParser.ConfigParser()
		self.config.optionxform = str
		self.user_home = os.path.expanduser('~')
		self.window.show()
		self.ini_check()

# Startup Checks
	def ini_check(self, data=None):
		ini_path = os.path.expanduser('~') + '/.config/dxf2emc'
		self.ini_file = ini_path + '/dxf2emc.ini'
		if not os.path.exists(ini_path):
			os.makedirs(ini_path, 0755)
		if not os.path.exists(self.ini_file):
			self.build_ini()
			message = 'Preferences File Created\nthis can be edited in the Preferences tab'
			result = self.ok_dialog(message)
			if result == gtk.RESPONSE_OK:
				print 'ok'
		if os.path.exists(self.ini_file):
			try:
				self.config.read(self.ini_file)
				self.on_revert_prefrences()
			except:
				message = 'The preferences file is corrupt.\nRebuild with default settings?'
				result = self.yesno_dialog(message)
				if result == gtk.RESPONSE_YES:
					self.build_ini()

	def build_ini(self):
		cfgfile = open(self.ini_file,'w+')
		self.config.add_section('Configuration')
		self.config.set('Configuration', 'TOLERANCE', '0.000001')
		self.config.set('Configuration', 'PATH', 'CCW')
		self.config.set('Configuration', 'SAVE_DIR', '')
		self.config.set('Configuration', 'OUTPUT', 'output.ngc')
		self.config.write(cfgfile)
		cfgfile.close

	def on_window_destroy(self, object, data=None):
		gtk.main_quit()

# File Menu Items
	def on_file_quit(self, menuitem, data=None):
		gtk.main_quit()

	def on_file_open(self, menuitem, data=None):
		self.fcd = gtk.FileChooserDialog("Open...", None,
			gtk.FILE_CHOOSER_ACTION_OPEN,
			(gtk.STOCK_CANCEL, gtk.RESPONSE_CANCEL, gtk.STOCK_OPEN, gtk.RESPONSE_OK))
		if len(self.current_folder) > 0:
			self.fcd.set_current_folder(self.current_folder)
		self.response = self.fcd.run()
		if self.response == gtk.RESPONSE_OK:
			self.status.set_text('File Selected %s' % self.fcd.get_filename())
			self.input_file = "-f=" + self.fcd.get_filename()
			self.current_folder = os.path.dirname(self.fcd.get_uri()[7:])
			self.file_analyze.set_sensitive(True)
			self.file_convert.set_sensitive(True)
			self.on_file_analyze()
		else:
			self.status.set_text('No File Open')
		self.fcd.destroy()

	def on_file_analyze(self, menuitem, data=None):
		print 'Analyze %s' % self.input_file
		command = "dxf2gcode -a"
		process = Popen(command, stdout=PIPE, stderr=PIPE, shell=True)
		output, error = process.communicate()

	def on_file_convert(self, input_file, data=None):
		if len(self.input_file) > 0:
			self.args = self.input_file
			self.result = Popen("dxf2gcode %s" %self.args, shell=True)
			if self.result == 0:
				self.status.set_text('Processing Complete')
			else:
				self.status.set_text('Error %d Processing %s' % (self.result, self.input_file))
		else:
			self.status.set_text('No File Open')

# View Menu Items
	def on_view_test(self, item, data=None):
		message = 'Do you like my test?\nPick one.'
		result = self.yesno_dialog(message)
		if result == gtk.RESPONSE_YES:
			print 'view yes'
		if result == gtk.RESPONSE_NO:
			print 'view no'
		if result == gtk.RESPONSE_DELETE_EVENT:
			print 'view delete'

# Help Menu Items
	def on_help_about(self, menuitem, data=None):
		self.response = self.aboutdialog.run()
		self.aboutdialog.hide()

# Setup Tab

# Preferences Tab
	def on_revert_prefrences(self, data=None):
		try:
			self.tolerance.set_text(self.config.get('Configuration', 'TOLERANCE'))
			self.path.set_text(self.config.get('Configuration', 'PATH'))
			self.save_directory.set_text(self.config.get('Configuration', 'SAVE_DIR'))
			self.output_name.set_text(self.config.get('Configuration', 'OUTPUT'))
		except:
			message = 'The preferences file is corrupt.\nRebuild with default settings?'
			result = yesno_dialog(message)
			if result == gtk.RESPONSE_YES:
				self.build_ini()

	def on_find_button(self, data=None):
		self.fcd = gtk.FileChooserDialog("Open...", None,
			gtk.FILE_CHOOSER_ACTION_SELECT_FOLDER,
			(gtk.STOCK_CANCEL, gtk.RESPONSE_CANCEL, gtk.STOCK_OPEN, gtk.RESPONSE_OK))
		self.fcd.set_current_folder(self.user_home)
		self.response = self.fcd.run()
		if self.response == gtk.RESPONSE_OK:
			print "Selected filepath: %s" % self.fcd.get_filename()
			self.save_directory.set_text(self.fcd.get_filename())
		self.fcd.destroy()

	def on_save_preferences(self, data=None):
		cfgfile = open(self.ini_file, 'w')
		self.config.set('Configuration', 'TOLERANCE', self.tolerance.get_text())
		self.config.set('Configuration', 'PATH', self.path.get_text())
		self.config.set('Configuration', 'SAVE_DIR', self.save_directory.get_text())
		self.config.set('Configuration', 'OUTPUT', self.output_name.get_text())
		self.config.write(cfgfile)
		cfgfile.close()

# Message Dialogs
	def yesno_dialog(self, message):
		md = gtk.MessageDialog(self.window,
		gtk.DIALOG_DESTROY_WITH_PARENT, gtk.MESSAGE_INFO,
		gtk.BUTTONS_YES_NO, message)
		result = md.run()
		md.destroy()
		return result

	def ok_dialog(self, message):
		md = gtk.MessageDialog(self.window,
		gtk.DIALOG_DESTROY_WITH_PARENT, gtk.MESSAGE_INFO,
		gtk.BUTTONS_OK, message)
		result = md.run()
		md.destroy()
		return result


if __name__ == '__main__':
  main = Buglump()
  gtk.main()
